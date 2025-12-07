package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"soda-interview/business/core/product"
	productstore "soda-interview/business/data/stores/product"
	tt "soda-interview/zarf/testing"
)

func Test_ProductService(t *testing.T) {
	c := tt.NewDBContainer(t)
	defer c.Teardown(t)
	c.Truncate(t)

	// Setup Store & Service
	pStore := productstore.NewStore(c.Log, c.DB)
	service := product.NewService(c.Log, pStore)
	ctx := context.Background()

	t.Run("Create_Success", func(t *testing.T) {
		np := product.NewProduct{
			Name:               "Soda",
			Description:        "Carbonated Drink",
			Price:              150,
			BuyerRewardPoints:  10,
			AuthorRewardPoints: 5,
		}

		p, err := service.Create(ctx, np)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if p.ID == "" {
			t.Error("Expected Product ID to be set")
		}
		if p.Name != np.Name {
			t.Errorf("Expected Name %s, got %s", np.Name, p.Name)
		}
		if p.Price != np.Price {
			t.Errorf("Expected Price %d, got %d", np.Price, p.Price)
		}
	})

	t.Run("Get_Success", func(t *testing.T) {
		// Create a product first
		np := product.NewProduct{
			Name:  "Get Me",
			Price: 200,
		}
		created, err := service.Create(ctx, np)
		if err != nil {
			t.Fatalf("Setup create failed: %v", err)
		}

		// Get it
		fetched, err := service.GetProduct(ctx, created.ID)
		if err != nil {
			t.Fatalf("GetProduct failed: %v", err)
		}

		if fetched.ID != created.ID {
			t.Errorf("Expected ID %s, got %s", created.ID, fetched.ID)
		}
		if fetched.Name != "Get Me" {
			t.Errorf("Expected Name 'Get Me', got %s", fetched.Name)
		}
	})

	t.Run("Get_Fail_NotFound", func(t *testing.T) {
		_, err := service.GetProduct(ctx, uuid.NewString())
		if err == nil {
			t.Fatal("Expected error for non-existent product, got nil")
		}
		if err.Error() == "" {
			t.Error("Expected error message, got empty")
		}
		// In a real app we might check for specific ErrNotFound wrapping
	})

	t.Run("List_Success", func(t *testing.T) {
		c.Truncate(t) // Clear DB for clean list count

		// Create 3 products
		for i := 0; i < 3; i++ {
			_, err := service.Create(ctx, product.NewProduct{Name: "Item", Price: 100})
			if err != nil {
				t.Fatalf("Setup create %d failed: %v", i, err)
			}
		}

		list, err := service.ListProducts(ctx)
		if err != nil {
			t.Fatalf("ListProducts failed: %v", err)
		}

		if len(list) != 3 {
			t.Errorf("Expected 3 products, got %d", len(list))
		}
	})
}
