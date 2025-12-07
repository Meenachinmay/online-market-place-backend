package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"soda-interview/business/data/stores/db"
	"soda-interview/foundation/logger"
)

var (
	ErrNotFound = errors.New("product not found")
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

// WithTx returns a new Store instance that uses the provided transaction.
func (s *Store) WithTx(tx pgx.Tx) *Store {
	return &Store{
		log: s.log,
		q:   s.q.WithTx(tx),
	}
}

func (s *Store) GetProduct(ctx context.Context, id string) (db.Product, error) {
	p, err := s.q.GetProduct(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return db.Product{}, ErrNotFound
		}
		return db.Product{}, fmt.Errorf("querying product: %w", err)
	}
	return p, nil
}

func (s *Store) ListProducts(ctx context.Context) ([]db.Product, error) {
	products, err := s.q.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	return products, nil
}

func (s *Store) CreateProduct(ctx context.Context, params db.CreateProductParams) (db.Product, error) {
	p, err := s.q.CreateProduct(ctx, params)
	if err != nil {
		return db.Product{}, fmt.Errorf("creating product: %w", err)
	}
	return p, nil
}