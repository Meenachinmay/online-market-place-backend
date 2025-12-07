package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"soda-interview/business/core/order"
	"soda-interview/business/core/referral-blog"
	"soda-interview/business/data/stores/db"
	orderstore "soda-interview/business/data/stores/order"
	productstore "soda-interview/business/data/stores/product"
	blogstore "soda-interview/business/data/stores/referral-blog"
	financestore "soda-interview/business/data/stores/soda-finance"
	tt "soda-interview/zarf/testing"
)

func Test_BlogToPurchaseFlow(t *testing.T) {
	c := tt.NewDBContainer(t)
	defer c.Teardown(t)
	c.Truncate(t)

	// Setup Stores
	oStore := orderstore.NewStore(c.Log, c.DB)
	pStore := productstore.NewStore(c.Log, c.DB)
	bStore := blogstore.NewStore(c.Log, c.DB)
	fStore := financestore.NewStore(c.Log, c.DB)

	// Setup Services
	blogService := referralblog.NewService(c.Log, bStore)
	orderService := order.NewService(c.Log, c.DB, oStore, pStore, bStore, fStore)

	ctx := context.Background()

	// Helper to create product (pre-requisite)
	createProduct := func(t *testing.T, price, buyerReward, authorReward int64) db.Product {
		id := uuid.NewString()
		p, err := pStore.CreateProduct(ctx, db.CreateProductParams{
			ID:                 id,
			Name:               "Flow Product",
			Description:        "A product for the flow test",
			Price:              price,
			BuyerRewardPoints:  int32(buyerReward),
			AuthorRewardPoints: int32(authorReward),
		})
		if err != nil {
			t.Fatalf("createProduct failed: %v", err)
		}
		return p
	}

	getWallet := func(t *testing.T, userID string) db.Wallet {
		w, err := fStore.GetWallet(ctx, userID)
		if err != nil {
			// Wallet might not exist yet if checks fail before creation
			return db.Wallet{}
		}
		return w
	}

	t.Run("Success_FullFlow", func(t *testing.T) {
		// 1. Setup
		authorID := uuid.NewString()
		buyerID := uuid.NewString()
		product := createProduct(t, 5000, 500, 250)

		// 2. Author creates Blog
		// Expectation: Author can create a blog linked to the product.
		newBlogReq := referralblog.NewBlog{
			AuthorID:  authorID,
			Content:   "This product changed my life! Buy it!",
			ProductID: product.ID,
		}

		blog, err := blogService.CreateBlog(ctx, newBlogReq)
		if err != nil {
			t.Fatalf("CreateBlog failed: %v", err)
		}

		if blog.ID == "" {
			t.Fatal("Expected Blog ID to be generated")
		}
		if blog.LinkedProductID != product.ID {
			t.Errorf("Expected linked product ID %s, got %s", product.ID, blog.LinkedProductID)
		}

		// 3. Buyer places Order via Blog
		// Expectation: Order is placed successfully, referencing the blog.
		placeOrderReq := order.PlaceOrderReq{
			BuyerID:   buyerID,
			ProductID: product.ID,
			BlogID:    blog.ID,
		}

		ord, err := orderService.PlaceOrder(ctx, placeOrderReq)
		if err != nil {
			t.Fatalf("PlaceOrder failed: %v", err)
		}

		if ord.Status != "CONFIRMED" {
			t.Errorf("Expected order status CONFIRMED, got %s", ord.Status)
		}

		// 4. Verify Points Awarded
		// Buyer Logic: First purchase of this product -> Gets BuyerRewardPoints (500)
		buyerWallet := getWallet(t, buyerID)
		if buyerWallet.SodaPoints != 500 {
			t.Errorf("Expected Buyer Points 500, got %d", buyerWallet.SodaPoints)
		}

		// Author Logic: Blog used for purchase -> Gets AuthorRewardPoints (250)
		authorWallet := getWallet(t, authorID)
		if authorWallet.SodaPoints != 250 {
			t.Errorf("Expected Author Points 250, got %d", authorWallet.SodaPoints)
		}
	})

	t.Run("Fail_CreateBlog_InvalidProduct", func(t *testing.T) {
		// Test that creating a blog with a non-existent product fails (if enforced by FK or service)
		// Since we have Foreign Key constraints in DB (usually), this should fail at DB level.
		
		authorID := uuid.NewString()
		nonExistentProductID := uuid.NewString()

		newBlogReq := referralblog.NewBlog{
			AuthorID:  authorID,
			Content:   "Fake product",
			ProductID: nonExistentProductID,
		}

		_, err := blogService.CreateBlog(ctx, newBlogReq)
		if err == nil {
			t.Fatal("Expected CreateBlog to fail with invalid product ID, but it succeeded")
		}
		// Optional: Check error message contains "foreign key constraint" or similar
	})
}
