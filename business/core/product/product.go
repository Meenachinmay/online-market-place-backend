package product

import (
	"context"
	"errors"
	"fmt"

	"soda-interview/business/data/stores/db"
	"soda-interview/foundation/logger"
)

var (
	ErrNotFound = errors.New("product not found")
)

type Product struct {
	ID                 string
	Name               string
	Description        string
	Price              int64
	BuyerRewardPoints  int32
	AuthorRewardPoints int32
}

type Service struct {
	log *logger.Logger
	store Storer
}

func NewService(log *logger.Logger, store Storer) *Service {
	return &Service{
		log: log,
		store: store,
	}
}

func (s *Service) GetProduct(ctx context.Context, id string) (Product, error) {
	p, err := s.store.GetProduct(ctx, id)
	if err != nil {
		// Check for not found error? 
		// The store returns db.Product and error.
		// If store implementation returns a specific error for Not Found, we can check it.
		// For now, assuming store returns error if not found.
		return Product{}, fmt.Errorf("querying product: %w", err)
	}
	// Ideally we handle "not found" mapping here if the store returns a specific error.
	return toProduct(p), nil
}

func (s *Service) ListProducts(ctx context.Context) ([]Product, error) {
	products, err := s.store.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	
	out := make([]Product, len(products))
	for i, p := range products {
		out[i] = toProduct(p)
	}
	return out, nil
}

func toProduct(dbP db.Product) Product {
	return Product{
		ID:                 dbP.ID,
		Name:               dbP.Name,
		Description:        dbP.Description,
		Price:              dbP.Price,
		BuyerRewardPoints:  dbP.BuyerRewardPoints,
		AuthorRewardPoints: dbP.AuthorRewardPoints,
	}
}