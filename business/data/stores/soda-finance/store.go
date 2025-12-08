package sodafinance

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
	ErrNotFound = errors.New("wallet not found")
	ErrInsufficientPoints = errors.New("insufficient points")
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

func (s *Store) GetWallet(ctx context.Context, userID string) (db.Wallet, error) {
	w, err := s.q.GetWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Wallet{}, ErrNotFound
		}
		return db.Wallet{}, fmt.Errorf("querying wallet: %w", err)
	}
	return w, nil
}

func (s *Store) GetOrCreateWallet(ctx context.Context, userID string) (db.Wallet, error) {
	w, err := s.q.CreateWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Wallet already exists, retrieve it
			return s.q.GetWallet(ctx, userID)
		}
		return db.Wallet{}, fmt.Errorf("creating wallet: %w", err)
	}
	return w, nil
}

func (s *Store) AddPoints(ctx context.Context, params db.AddPointsParams) (db.Wallet, error) {
	w, err := s.q.AddPoints(ctx, params)
	if err != nil {
		return db.Wallet{}, fmt.Errorf("adding points: %w", err)
	}
	return w, nil
}

func (s *Store) AddBalance(ctx context.Context, params db.AddBalanceParams) (db.Wallet, error) {
	w, err := s.q.AddBalance(ctx, params)
	if err != nil {
		return db.Wallet{}, fmt.Errorf("adding balance: %w", err)
	}
	return w, nil
}

func (s *Store) ConvertPointsToBalance(ctx context.Context, params db.ConvertPointsToBalanceParams) (db.Wallet, error) {
	w, err := s.q.ConvertPointsToBalance(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Wallet{}, ErrInsufficientPoints
		}
		return db.Wallet{}, fmt.Errorf("converting points: %w", err)
	}
	return w, nil
}

func (s *Store) CreateTransaction(ctx context.Context, params db.CreateTransactionParams) (db.Transaction, error) {
	t, err := s.q.CreateTransaction(ctx, params)
	if err != nil {
		return db.Transaction{}, fmt.Errorf("creating transaction: %w", err)
	}
	return t, nil
}