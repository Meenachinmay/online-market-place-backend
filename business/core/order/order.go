package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"soda-interview/business/data/stores/db"
	orderstore "soda-interview/business/data/stores/order"
	productstore "soda-interview/business/data/stores/product"
	blogstore "soda-interview/business/data/stores/referral-blog"
	financestore "soda-interview/business/data/stores/soda-finance"
	"soda-interview/foundation/logger"
)

type Order struct {
	ID        string
	BuyerID   string
	ProductID string
	BlogID    string
	Amount    int64
	Status    string
	CreatedAt int64
}

type PlaceOrderReq struct {
	BuyerID   string
	ProductID string
	BlogID    string
}

type Service struct {
	log          *logger.Logger
	pool         *pgxpool.Pool
	orderStore   *orderstore.Store
	productStore *productstore.Store
	blogStore    *blogstore.Store
	financeStore *financestore.Store
}

func NewService(
	log *logger.Logger,
	pool *pgxpool.Pool,
	orderStore *orderstore.Store,
	productStore *productstore.Store,
	blogStore *blogstore.Store,
	financeStore *financestore.Store,
) *Service {
	return &Service{
		log:          log,
		pool:         pool,
		orderStore:   orderStore,
		productStore: productStore,
		blogStore:    blogStore,
		financeStore: financeStore,
	}
}

func (s *Service) PlaceOrder(ctx context.Context, req PlaceOrderReq) (Order, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Order{}, fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Create transactional stores
	qTxOrder := s.orderStore.WithTx(tx)
	qTxProduct := s.productStore.WithTx(tx)
	qTxBlog := s.blogStore.WithTx(tx)
	qTxFinance := s.financeStore.WithTx(tx)

	product, err := qTxProduct.GetProduct(ctx, req.ProductID)
	if err != nil {
		return Order{}, fmt.Errorf("getting product: %w", err)
	}

	blog, err := qTxBlog.GetBlog(ctx, req.BlogID)
	if err != nil {
		return Order{}, fmt.Errorf("getting blog: %w", err)
	}

	count, err := qTxOrder.CountOrdersByBuyer(ctx, req.BuyerID)
	if err != nil {
		return Order{}, fmt.Errorf("counting orders: %w", err)
	}

	isFirstPurchase := count == 0

	orderID := uuid.NewString()
	now := time.Now()
	dbOrder, err := qTxOrder.CreateOrder(ctx, db.CreateOrderParams{
		ID:        orderID,
		BuyerID:   req.BuyerID,
		ProductID: req.ProductID,
		BlogID:    req.BlogID,
		Amount:    product.Price,
		Status:    "CONFIRMED",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return Order{}, fmt.Errorf("creating order: %w", err)
	}

	if _, err := qTxFinance.CreateWallet(ctx, req.BuyerID); err != nil {
		return Order{}, fmt.Errorf("ensuring buyer wallet: %w", err)
	}

	if isFirstPurchase {
		buyerPoints := int64(product.BuyerRewardPoints)
		if buyerPoints > 0 {
			if _, err := qTxFinance.AddPoints(ctx, db.AddPointsParams{
				Amount: buyerPoints,
				UserID: req.BuyerID,
			}); err != nil {
				return Order{}, fmt.Errorf("adding buyer points: %w", err)
			}

			if _, err := qTxFinance.CreateTransaction(ctx, db.CreateTransactionParams{
				ID:             uuid.NewString(),
				UserID:         req.BuyerID,
				Type:           "EARNED",
				Amount:         buyerPoints,
				RelatedOrderID: pgtype.Text{String: orderID, Valid: true},
			}); err != nil {
				return Order{}, fmt.Errorf("logging buyer transaction: %w", err)
			}
		}
	}

	authorID := blog.AuthorID
	if _, err := qTxFinance.CreateWallet(ctx, authorID); err != nil {
		return Order{}, fmt.Errorf("ensuring author wallet: %w", err)
	}

	authorPoints := int64(product.AuthorRewardPoints)
	if authorPoints > 0 {
		if _, err := qTxFinance.AddPoints(ctx, db.AddPointsParams{
			Amount: authorPoints,
			UserID: authorID,
		}); err != nil {
			return Order{}, fmt.Errorf("adding author points: %w", err)
		}

		if _, err := qTxFinance.CreateTransaction(ctx, db.CreateTransactionParams{
			ID:             uuid.NewString(),
			UserID:         authorID,
			Type:           "EARNED",
			Amount:         authorPoints,
			RelatedOrderID: pgtype.Text{String: orderID, Valid: true},
		}); err != nil {
			return Order{}, fmt.Errorf("logging author transaction: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Order{}, fmt.Errorf("committing transaction: %w", err)
	}

	return Order{
		ID:        dbOrder.ID,
		BuyerID:   dbOrder.BuyerID,
		ProductID: dbOrder.ProductID,
		BlogID:    dbOrder.BlogID,
		Amount:    dbOrder.Amount,
		Status:    dbOrder.Status,
		CreatedAt: dbOrder.CreatedAt.Time.Unix(),
	}, nil
}
