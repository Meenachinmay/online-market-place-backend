package finance

import (
	"context"
	"github.com/jackc/pgx/v5"
	"soda-interview/business/data/stores/db"
)

type Storer interface {
	GetWallet(ctx context.Context, userID string) (db.Wallet, error)
	CreateWallet(ctx context.Context, userID string) (db.Wallet, error)
	AddPoints(ctx context.Context, params db.AddPointsParams) (db.Wallet, error)
	AddBalance(ctx context.Context, params db.AddBalanceParams) (db.Wallet, error)
	ConvertPointsToBalance(ctx context.Context, params db.ConvertPointsToBalanceParams) (db.Wallet, error)
	CreateTransaction(ctx context.Context, params db.CreateTransactionParams) (db.Transaction, error)
	
	// WithTx returns a version of the store that runs in the transaction.
	// For interface compliance, we might return 'Storer' or use a specific mechanism.
	// To allow 'WithTx' in interface, the return type must be 'Storer'.
	// But implementation returns *Store. This works in Go if *Store implements Storer.
	WithTx(tx pgx.Tx) Storer
}
