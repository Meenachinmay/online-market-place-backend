package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"soda-interview/business/core/order"
	"soda-interview/business/data/stores/db"
	orderstore "soda-interview/business/data/stores/order"
	productstore "soda-interview/business/data/stores/product"
	blogstore "soda-interview/business/data/stores/referral-blog"
	financestore "soda-interview/business/data/stores/soda-finance"
	tt "soda-interview/zarf/testing"
)

func Test_PlaceOrder(t *testing.T) {
	c := tt.NewDBContainer(t)
	defer c.Teardown(t)
	c.Truncate(t)

	// Setup Stores
	oStore := orderstore.NewStore(c.Log, c.DB)
	pStore := productstore.NewStore(c.Log, c.DB)
	bStore := blogstore.NewStore(c.Log, c.DB)
	fStore := financestore.NewStore(c.Log, c.DB)

	// Setup Service
	service := order.NewService(c.Log, c.DB, oStore, pStore, bStore, fStore)
	ctx := context.Background()

	// Helpers
	createProduct := func(t *testing.T, price, buyerReward, authorReward int64) db.Product {
		id := uuid.NewString()
		p, err := pStore.CreateProduct(ctx, db.CreateProductParams{
			ID:                 id,
			Name:               "Test Product",
			Description:        "Desc",
			Price:              price,
			BuyerRewardPoints:  int32(buyerReward),
			AuthorRewardPoints: int32(authorReward),
		})
		if err != nil {
			t.Fatalf("createProduct failed: %v", err)
		}
		return p
	}

	createBlog := func(t *testing.T, authorID, productID string) db.Blog {
		id := uuid.NewString()
		b, err := bStore.CreateBlog(ctx, db.CreateBlogParams{
			ID:        id,
			AuthorID:  authorID,
			Content:   "Check this out!",
			ProductID: productID,
		})
		if err != nil {
			t.Fatalf("createBlog failed: %v", err)
		}
		return b
	}

	getWallet := func(t *testing.T, userID string) db.Wallet {
		w, err := fStore.GetWallet(ctx, userID)
		if err != nil {
			t.Fatalf("getWallet failed: %v", err)
		}
		return w
	}

	t.Run("Success_FirstTimePurchase", func(t *testing.T) {
		c.Truncate(t)
		buyerID := uuid.NewString()
		authorID := uuid.NewString()

		prod := createProduct(t, 1000, 100, 50)
		blog := createBlog(t, authorID, prod.ID)

		req := order.PlaceOrderReq{
			BuyerID:   buyerID,
			ProductID: prod.ID,
			BlogID:    blog.ID,
		}

		// Execute
		ord, err := service.PlaceOrder(ctx, req)
		if err != nil {
			t.Fatalf("PlaceOrder failed: %v", err)
		}
		if ord.ID == "" {
			t.Error("expected Order ID to be set")
		}
		if ord.BuyerID != buyerID {
			t.Errorf("expected BuyerID %s, got %s", buyerID, ord.BuyerID)
		}
		if ord.Amount != 1000 {
			t.Errorf("expected Amount 1000, got %d", ord.Amount)
		}

		// Verify Buyer Wallet (Should have 100 points)
		buyerWallet := getWallet(t, buyerID)
		if buyerWallet.SodaPoints != 100 {
			t.Errorf("expected Buyer Points 100, got %d", buyerWallet.SodaPoints)
		}

		// Verify Author Wallet (Should have 50 points)
		authorWallet := getWallet(t, authorID)
		if authorWallet.SodaPoints != 50 {
			t.Errorf("expected Author Points 50, got %d", authorWallet.SodaPoints)
		}
	})

	t.Run("Success_SecondTimePurchase_SameProduct", func(t *testing.T) {
		c.Truncate(t)
		buyerID := uuid.NewString()
		authorID := uuid.NewString()

		prod := createProduct(t, 1000, 100, 50)
		blog := createBlog(t, authorID, prod.ID)

		req := order.PlaceOrderReq{
			BuyerID:   buyerID,
			ProductID: prod.ID,
			BlogID:    blog.ID,
		}

		// First Order
		_, err := service.PlaceOrder(ctx, req)
		if err != nil {
			t.Fatalf("First PlaceOrder failed: %v", err)
		}

		// Second Order
		time.Sleep(10 * time.Millisecond)
		ord2, err := service.PlaceOrder(ctx, req)
		if err != nil {
			t.Fatalf("Second PlaceOrder failed: %v", err)
		}
		if ord2.ID == "" {
			t.Error("expected Order ID to be set for second order")
		}

		// Verify Buyer Wallet (Should STILL have 100 points - no double dip)
		buyerWallet := getWallet(t, buyerID)
		if buyerWallet.SodaPoints != 100 {
			t.Errorf("expected Buyer Points 100 (no double dip), got %d", buyerWallet.SodaPoints)
		}

		// Verify Author Wallet (Should have 100 points - 50 + 50)
		authorWallet := getWallet(t, authorID)
		if authorWallet.SodaPoints != 100 {
			t.Errorf("expected Author Points 100, got %d", authorWallet.SodaPoints)
		}
	})

	t.Run("Success_DifferentProduct_SameBuyer", func(t *testing.T) {
		c.Truncate(t)
		buyerID := uuid.NewString()
		authorID := uuid.NewString()

		// Product 1
		prod1 := createProduct(t, 1000, 100, 50)
		blog1 := createBlog(t, authorID, prod1.ID)

		// Product 2
		prod2 := createProduct(t, 2000, 200, 50)
		blog2 := createBlog(t, authorID, prod2.ID)

		// Buy Product 1
		_, err := service.PlaceOrder(ctx, order.PlaceOrderReq{
			BuyerID:   buyerID,
			ProductID: prod1.ID,
			BlogID:    blog1.ID,
		})
		if err != nil {
			t.Fatalf("PlaceOrder Product 1 failed: %v", err)
		}

		// Buy Product 2
		_, err = service.PlaceOrder(ctx, order.PlaceOrderReq{
			BuyerID:   buyerID,
			ProductID: prod2.ID,
			BlogID:    blog2.ID,
		})
		if err != nil {
			t.Fatalf("PlaceOrder Product 2 failed: %v", err)
		}

		// Verify Buyer Wallet 
		// Expectation: 100 (Prod1) + 200 (Prod2) = 300
		buyerWallet := getWallet(t, buyerID)
		if buyerWallet.SodaPoints != 300 {
			t.Errorf("expected Buyer Points 300 (100+200), got %d. NOTE: If this is 100, the 'First Purchase' logic is flawed.", buyerWallet.SodaPoints)
		}
	})
}
