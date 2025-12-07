package product

import (
	"context"

	"soda-interview/business/core/product"
	productv1 "soda-interview/foundation/proto/product/v1"
)

type Handler struct {
	productv1.UnimplementedProductServiceServer
	Service *product.Service
}

func (h *Handler) GetProduct(ctx context.Context, req *productv1.ProductRequest) (*productv1.Product, error) {
	p, err := h.Service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &productv1.Product{
		Id:                p.ID,
		Name:              p.Name,
		Description:       p.Description,
		Price:             p.Price,
		BuyerRewardPoints: p.BuyerRewardPoints,
	}, nil
}

func (h *Handler) ListProducts(ctx context.Context, _ *productv1.Empty) (*productv1.ProductList, error) {
	products, err := h.Service.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*productv1.Product, len(products))
	for i, p := range products {
		list[i] = &productv1.Product{
			Id:                p.ID,
			Name:              p.Name,
			Description:       p.Description,
			Price:             p.Price,
			BuyerRewardPoints: p.BuyerRewardPoints,
		}
	}

	return &productv1.ProductList{Products: list}, nil
}