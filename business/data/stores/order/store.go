package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"soda-interview/business/data/stores/db"
	"soda-interview/foundation/logger"
)

type Store struct {
	log *logger.Logger
	q   *db.Queries
}

func NewStore(log *logger.Logger, pool *pgxpool.Pool) *Store {
	return &Store{
		log: log,
		q:   db.New(pool),
	}
}

func (s *Store) WithTx(tx pgx.Tx) *Store {
	return &Store{
		log: s.log,
		q:   s.q.WithTx(tx),
	}
}

func (s *Store) CreateOrder(ctx context.Context, params db.CreateOrderParams) (db.Order, error) {
	o, err := s.q.CreateOrder(ctx, params)
	if err != nil {
		return db.Order{}, fmt.Errorf("creating order: %w", err)
	}
	return o, nil
}

func (s *Store) CountOrdersByBuyer(ctx context.Context, buyerID string) (int64, error) {
	c, err := s.q.CountOrdersByBuyer(ctx, buyerID)
	if err != nil {
		return 0, fmt.Errorf("counting orders: %w", err)
	}
	return c, nil
}

func (s *Store) CountOrdersByBuyerAndProduct(ctx context.Context, buyerID, productID string) (int64, error) {
	c, err := s.q.CountOrdersByBuyerAndProduct(ctx, db.CountOrdersByBuyerAndProductParams{
		BuyerID:   buyerID,
		ProductID: productID,
	})
	if err != nil {
		return 0, fmt.Errorf("counting orders by product: %w", err)
	}
	return c, nil
}