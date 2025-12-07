package product

import (
	"context"

	"soda-interview/business/data/stores/db"
)

// Storer defines the behavior required by the product service.
type Storer interface {
	CreateProduct(ctx context.Context, params db.CreateProductParams) (db.Product, error)
	GetProduct(ctx context.Context, id string) (db.Product, error)
	ListProducts(ctx context.Context) ([]db.Product, error)
}
