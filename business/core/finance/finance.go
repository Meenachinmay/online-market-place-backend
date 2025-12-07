package finance

import (
	"context"
	"errors"
	"fmt"

	"soda-interview/business/data/stores/db"
	sodafinance "soda-interview/business/data/stores/soda-finance"
	"soda-interview/foundation/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("wallet not found")
	ErrInsufficientPoints = errors.New("insufficient points")
)

type Wallet struct {
	UserID      string
	SodaPoints  int64
	SodaBalance int64
}

type Service struct {
	log *logger.Logger
	pool *pgxpool.Pool
	store *sodafinance.Store
}

func NewService(log *logger.Logger, pool *pgxpool.Pool, store *sodafinance.Store) *Service {
	return &Service{
		log:  log,
		pool: pool,
		store: store,
	}
}

func (s *Service) GetWallet(ctx context.Context, userID string) (Wallet, error) {
	w, err := s.store.GetWallet(ctx, userID)
	if err != nil {
		// handle specific store errors if needed
		return Wallet{}, fmt.Errorf("querying wallet: %w", err)
	}
	return toWallet(w), nil
}

func (s *Service) EnsureWalletExists(ctx context.Context, userID string) error {
	_, err := s.store.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return fmt.Errorf("creating wallet: %w", err)
	}
	return nil
}

func (s *Service) ConvertPoints(ctx context.Context, userID string, pointsToConvert int64) (Wallet, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Wallet{}, fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Use the transactional store
	txStore := s.store.WithTx(tx)

	w, err := txStore.GetWallet(ctx, userID)
	if err != nil {
		return Wallet{}, fmt.Errorf("getting wallet: %w", err)
	}

	if w.SodaPoints <= 1000 {
		return Wallet{}, fmt.Errorf("%w: must have > 1000 points, have %d", ErrInsufficientPoints, w.SodaPoints)
	}

	amount := pointsToConvert
	if amount <= 0 {
		amount = w.SodaPoints
	}
	
	if amount > w.SodaPoints {
		return Wallet{}, fmt.Errorf("%w: requesting %d, have %d", ErrInsufficientPoints, amount, w.SodaPoints)
	}
	
	yen := amount / 2
	pointsDeducted := yen * 2
	
	if pointsDeducted == 0 {
		if err := tx.Commit(ctx); err != nil {
			return Wallet{}, fmt.Errorf("commit: %w", err)
		}
		return toWallet(w), nil
	}

	updatedW, err := txStore.ConvertPointsToBalance(ctx, db.ConvertPointsToBalanceParams{
		PointsDeducted: pointsDeducted,
		BalanceAdded:   yen,
		UserID:         userID,
	})
	if err != nil {
		return Wallet{}, fmt.Errorf("converting points: %w", err)
	}

	_, err = txStore.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:             uuid.NewString(),
		UserID:         userID,
		Type:           "CONVERTED",
		Amount:         pointsDeducted,
		RelatedOrderID: pgtype.Text{Valid: false},
	})
	if err != nil {
		return Wallet{}, fmt.Errorf("creating transaction log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Wallet{}, fmt.Errorf("committing transaction: %w", err)
	}
	
	return toWallet(updatedW), nil
}

func toWallet(w db.Wallet) Wallet {
	return Wallet{
		UserID:      w.UserID,
		SodaPoints:  w.SodaPoints,
		SodaBalance: w.SodaBalance,
	}
}
