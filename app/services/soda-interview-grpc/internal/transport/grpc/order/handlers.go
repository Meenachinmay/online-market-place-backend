package order

import (
	"context"

	"soda-interview/business/core/order"
	orderv1 "soda-interview/foundation/proto/order/v1"
)

type Handler struct {
	orderv1.UnimplementedOrderServiceServer
	Service *order.Service
}

func (h *Handler) PlaceOrder(ctx context.Context, req *orderv1.PlaceOrderRequest) (*orderv1.OrderResponse, error) {
	o, err := h.Service.PlaceOrder(ctx, order.PlaceOrderReq{
		BuyerID:   req.BuyerId,
		ProductID: req.ProductId,
		BlogID:    req.BlogId,
	})
	if err != nil {
		return nil, err
	}

	return &orderv1.OrderResponse{
		Order: &orderv1.Order{
			Id:        o.ID,
			BuyerId:   o.BuyerID,
			ProductId: o.ProductID,
			Amount:    o.Amount,
			Status:    o.Status,
			CreatedAt: o.CreatedAt,
		},
	}, nil
}