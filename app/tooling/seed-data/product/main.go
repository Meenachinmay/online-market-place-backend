package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"soda-interview/business/core/product"
	"soda-interview/business/data/stores/db"
	productstore "soda-interview/business/data/stores/product"
	"soda-interview/foundation/config"
	"soda-interview/foundation/database/postgres"
	"soda-interview/foundation/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	log := logger.New(os.Stdout, "INFO")

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../../../../")
	configPath := filepath.Join(projectRoot, "foundation/config")

	cfg, err := config.LoadWithPath(configPath, "local")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbPool, err := postgres.New(ctx, cfg.GetDatabaseDSN())
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer dbPool.Close()

	pStore := productstore.NewStore(log, dbPool)
	_ = product.NewService(log, pStore) // Ensure service matches if needed, but we use store for explicit ID

	products := []struct {
		ID                 string
		Name               string
		Desc               string
		Price              int64
		BuyerRewardPoints  int32
		AuthorRewardPoints int32
	}{
		{"68d4d95a-5df3-40c9-b0e2-fb5c30904e09", "Soda Air Max 1", "Classic comfort and style for everyday wear.", 100, 10, 5},
		{"07ea4a41-ea43-465b-a121-cedc5becf52e", "Soda Runner Low", "Lightweight running shoes for speed.", 200, 20, 10},
		{"a3a35d83-fe83-45a1-b6f1-fb6f03debf84", "Soda High Tops", "Vintage style high tops with modern grip.", 300, 30, 15},
		{"719b2afe-034d-4ec2-8b51-aae6fced2884", "Soda Court Vision", "Basketball inspired sneakers for the court.", 400, 40, 20},
		{"126f3a3e-302d-4286-8971-c3e98f5778ef", "Soda Street King", "Urban street style with premium leather.", 500, 50, 25},
		{"d142a19c-6e05-4183-9c48-99330d82d54d", "Soda Canvas Slip-On", "Casual slip-ons for a relaxed vibe.", 600, 60, 30},
		{"f1fa509f-57a5-4d28-a111-863d470978da", "Soda Trail Blazer", "Rugged outsole for off-road adventures.", 700, 70, 35},
		{"25d45deb-c5ab-449d-b9a8-5f01d061972e", "Soda Retro 90s", "Throwback design with chunky soles.", 800, 80, 40},
		{"54972b8c-b1fd-4cab-8788-53a5fd193a1f", "Soda Knit Runner", "Breathable knit upper for maximum airflow.", 900, 90, 45},
		{"61aad777-8e7c-4135-ab2b-590561bcaea5", "Soda Pro Skater", "Durable suede reinforced for skating.", 1000, 100, 50},
		{"2db45c6f-e98a-4d75-a3de-bd39a4bd4542", "Soda Elite Racer", "Professional grade marathon shoes.", 1100, 110, 55},
		{"4852788b-47eb-4ebf-8699-13b8ad58ff76", "Soda Limited Edition", "Exclusive colorway, limited stock.", 1200, 120, 60},
		{"04777eaa-7e1a-4418-88ae-1eda740b3c3d", "Soda Tech Future", "Futuristic design with auto-lacing tech.", 1300, 130, 65},
		{"838d7318-5ea5-444e-bfd1-bec8e362617e", "Soda Classic White", "Minimalist white sneakers that go with everything.", 1400, 140, 70},
		{"257270b6-f9e1-40af-80d2-e11980800e58", "Soda Midnight Black", "Sleek all-black design for night outs.", 1500, 150, 75},
	}

	log.InfoContext(ctx, "Starting product seed...")
	for _, p := range products {
		_, err := pStore.GetProduct(ctx, p.ID)
		if err == nil {
			log.InfoContext(ctx, "Product already exists, skipping", "id", p.ID)
			continue
		}

		_, err = pStore.CreateProduct(ctx, db.CreateProductParams{
			ID:                 p.ID,
			Name:               p.Name,
			Description:        p.Desc,
			Price:              p.Price,
			BuyerRewardPoints:  p.BuyerRewardPoints,
			AuthorRewardPoints: p.AuthorRewardPoints,
		})
		if err != nil {
			return fmt.Errorf("creating product %s: %w", p.ID, err)
		}
		log.InfoContext(ctx, "Created product", "id", p.ID, "name", p.Name)
	}
	
	log.InfoContext(ctx, "Product seed completed successfully")
	return nil
}
