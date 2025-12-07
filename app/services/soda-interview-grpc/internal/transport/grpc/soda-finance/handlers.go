package sodafinance

import (
	"context"

	"soda-interview/business/core/finance"
	financev1 "soda-interview/foundation/proto/soda-finance/v1"
)

type Handler struct {
	financev1.UnimplementedFinanceServiceServer
	Service *finance.Service
}

func (h *Handler) GetWallet(ctx context.Context, req *financev1.UserRequest) (*financev1.Wallet, error) {
	w, err := h.Service.GetWallet(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &financev1.Wallet{
		UserId:      w.UserID,
		SodaPoints:  w.SodaPoints,
		SodaBalance: w.SodaBalance,
	}, nil
}

func (h *Handler) ConvertPoints(ctx context.Context, req *financev1.ConvertRequest) (*financev1.Wallet, error) {
	w, err := h.Service.ConvertPoints(ctx, req.UserId, req.PointsToConvert)
	if err != nil {
		return nil, err
	}

	return &financev1.Wallet{
		UserId:      w.UserID,
		SodaPoints:  w.SodaPoints,
		SodaBalance: w.SodaBalance,
	}, nil
}