package order

import (
	"context"
	"fmt"
	"time"

	"soda-interview/business/data/stores/db"
	orderstore "soda-interview/business/data/stores/order"
	productstore "soda-interview/business/data/stores/product"
	blogstore "soda-interview/business/data/stores/referral-blog"
	financestore "soda-interview/business/data/stores/soda-finance"
	"soda-interview/foundation/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

	count, err := qTxOrder.CountOrdersByBuyerAndProduct(ctx, req.BuyerID, req.ProductID)
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

	if _, err := qTxFinance.GetOrCreateWallet(ctx, req.BuyerID); err != nil {
		return Order{}, fmt.Errorf("ensuring buyer wallet: %w", err)
	}

	if isFirstPurchase {
		if err := s.distributeBuyerRewards(ctx, qTxFinance, req.BuyerID, product.BuyerRewardPoints, orderID); err != nil {
			return Order{}, fmt.Errorf("distributing buyer rewards: %w", err)
		}
	}

	authorID := blog.AuthorID
	if _, err := qTxFinance.GetOrCreateWallet(ctx, authorID); err != nil {
		return Order{}, fmt.Errorf("ensuring author wallet: %w", err)
	}

	if err := s.distributeAuthorRewards(ctx, qTxFinance, authorID, product.AuthorRewardPoints, orderID); err != nil {
		return Order{}, fmt.Errorf("distributing author rewards: %w", err)
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

func (s *Service) distributeBuyerRewards(ctx context.Context, txFinance *financestore.Store, buyerID string, points int32, orderID string) error {
	amount := int64(points)
	if amount <= 0 {
		return nil
	}

	if _, err := txFinance.AddPoints(ctx, db.AddPointsParams{
		Amount: amount,
		UserID: buyerID,
	}); err != nil {
		return fmt.Errorf("adding points: %w", err)
	}

	if _, err := txFinance.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:             uuid.NewString(),
		UserID:         buyerID,
		Type:           "EARNED",
		Amount:         amount,
		RelatedOrderID: pgtype.Text{String: orderID, Valid: true},
	}); err != nil {
		return fmt.Errorf("logging transaction: %w", err)
	}
	return nil
}

func (s *Service) distributeAuthorRewards(ctx context.Context, txFinance *financestore.Store, authorID string, points int32, orderID string) error {
	amount := int64(points)
	if amount <= 0 {
		return nil
	}

	if _, err := txFinance.AddPoints(ctx, db.AddPointsParams{
		Amount: amount,
		UserID: authorID,
	}); err != nil {
		return fmt.Errorf("adding points: %w", err)
	}

	if _, err := txFinance.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:             uuid.NewString(),
		UserID:         authorID,
		Type:           "EARNED",
		Amount:         amount,
		RelatedOrderID: pgtype.Text{String: orderID, Valid: true},
	}); err != nil {
		return fmt.Errorf("logging transaction: %w", err)
	}
	return nil
}
