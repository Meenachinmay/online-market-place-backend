package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"soda-interview/business/core/finance"
	"soda-interview/business/data/stores/db"
	financestore "soda-interview/business/data/stores/soda-finance"
	tt "soda-interview/zarf/testing"
)

func Test_ConvertPoints(t *testing.T) {
	c := tt.NewDBContainer(t)
	defer c.Teardown(t)
	c.Truncate(t)

	// Setup Stores
	fStore := financestore.NewStore(c.Log, c.DB)

	// Setup Service
	service := finance.NewService(c.Log, c.DB, fStore)
	ctx := context.Background()

	// Helpers
	setupWallet := func(t *testing.T, points int64) string {
		userID := uuid.NewString()
		_, err := fStore.GetOrCreateWallet(ctx, userID)
		if err != nil {
			t.Fatalf("setupWallet: create failed: %v", err)
		}
		if points > 0 {
			_, err = fStore.AddPoints(ctx, db.AddPointsParams{
				Amount: points,
				UserID: userID,
			})
			if err != nil {
				t.Fatalf("setupWallet: add points failed: %v", err)
			}
		}
		return userID
	}

	getWallet := func(t *testing.T, userID string) db.Wallet {
		w, err := fStore.GetWallet(ctx, userID)
		if err != nil {
			t.Fatalf("getWallet failed: %v", err)
		}
		return w
	}

	t.Run("Fail_InsufficientPoints_UnderThreshold", func(t *testing.T) {
		userID := setupWallet(t, 500) // Less than 1000
		
		_, err := service.ConvertPoints(ctx, userID, 500)
		if err == nil {
			t.Fatal("expected error for < 1000 points, got nil")
		}
	})

	t.Run("Fail_RequestMoreThanAvailable", func(t *testing.T) {
		userID := setupWallet(t, 2000)
		
		_, err := service.ConvertPoints(ctx, userID, 3000)
		if err == nil {
			t.Fatal("expected error for requesting more than available, got nil")
		}
	})

	t.Run("Success_ConvertSpecificAmount", func(t *testing.T) {
		userID := setupWallet(t, 2000)

		// Convert 1000 points. 2 points = 1 yen. Expect 500 yen balance.
		w, err := service.ConvertPoints(ctx, userID, 1000)
		if err != nil {
			t.Fatalf("ConvertPoints failed: %v", err)
		}

		if w.SodaPoints != 1000 { // 2000 - 1000
			t.Errorf("expected 1000 points remaining, got %d", w.SodaPoints)
		}
		if w.SodaBalance != 500 { // 1000 / 2
			t.Errorf("expected 500 balance, got %d", w.SodaBalance)
		}

		// Verify DB state
		dbW := getWallet(t, userID)
		if dbW.SodaPoints != 1000 || dbW.SodaBalance != 500 {
			t.Errorf("db mismatch: %+v", dbW)
		}
	})

	t.Run("Success_ConvertAll_Default", func(t *testing.T) {
		userID := setupWallet(t, 2000)

		// Convert 0 => convert all.
		w, err := service.ConvertPoints(ctx, userID, 0)
		if err != nil {
			t.Fatalf("ConvertPoints failed: %v", err)
		}

		if w.SodaPoints != 0 {
			t.Errorf("expected 0 points remaining, got %d", w.SodaPoints)
		}
		if w.SodaBalance != 1000 { // 2000 / 2
			t.Errorf("expected 1000 balance, got %d", w.SodaBalance)
		}
	})

	t.Run("Success_OddPoints_Rounding", func(t *testing.T) {
		// If we have 1001 points. 
		// Convert 1001. 
		// Yen = 1001 / 2 = 500. 
		// Points Deducted = 500 * 2 = 1000.
		// Remaining = 1.
		userID := setupWallet(t, 1001)

		w, err := service.ConvertPoints(ctx, userID, 0)
		if err != nil {
			t.Fatalf("ConvertPoints failed: %v", err)
		}

		if w.SodaPoints != 1 {
			t.Errorf("expected 1 point remaining, got %d", w.SodaPoints)
		}
		if w.SodaBalance != 500 {
			t.Errorf("expected 500 balance, got %d", w.SodaBalance)
		}
	})
}
