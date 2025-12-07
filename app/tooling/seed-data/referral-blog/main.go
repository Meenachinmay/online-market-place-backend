package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"soda-interview/business/core/referral-blog"
	"soda-interview/business/data/stores/db"
	blogstore "soda-interview/business/data/stores/referral-blog"
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

	bStore := blogstore.NewStore(log, dbPool)
	_ = referralblog.NewService(log, bStore)

	// Hardcoded data matching the SQL generation plan for consistency
	blogs := []struct {
		ID        string
		AuthorID  string
		Content   string
		ProductID string
	}{
		{"c0fa0d06-156a-4ddf-bfe1-e048b5597dba", "91183969-4b3a-4692-8ec6-74e2caf131af", "These Soda Air Max 1s are incredibly comfortable for walking all day. Highly recommend!", "68d4d95a-5df3-40c9-b0e2-fb5c30904e09"},
		{"cbd832e2-2957-4492-b7d7-b63e8c31711a", "91183969-4b3a-4692-8ec6-74e2caf131af", "Just got my Soda Runner Lows. They are so light, it feels like running on clouds.", "07ea4a41-ea43-465b-a121-cedc5becf52e"},
		{"610c6ca6-ec21-47e9-8831-b78cc209c124", "91183969-4b3a-4692-8ec6-74e2caf131af", "The vintage look on these Soda High Tops is fire. Great grip too.", "a3a35d83-fe83-45a1-b6f1-fb6f03debf84"},
		{"d4bc64a4-26de-464f-8cb1-df29ac3c0355", "91183969-4b3a-4692-8ec6-74e2caf131af", "Soda Court Vision has improved my game. Excellent ankle support.", "719b2afe-034d-4ec2-8b51-aae6fced2884"},
		{"be6714ee-49f9-45cb-bce4-0e8bea8b14ce", "91183969-4b3a-4692-8ec6-74e2caf131af", "Rocking the Soda Street Kings today. The leather quality is premium.", "126f3a3e-302d-4286-8971-c3e98f5778ef"},
		{"82ce1b67-f63a-47dd-b105-79243a4e9774", "91183969-4b3a-4692-8ec6-74e2caf131af", "Soda Canvas Slip-Ons are my go-to for quick errands. Super easy.", "d142a19c-6e05-4183-9c48-99330d82d54d"},
		{"76f3f15a-de44-4231-9b53-a6e1d4ec542b", "91183969-4b3a-4692-8ec6-74e2caf131af", "Took the Soda Trail Blazers hiking this weekend. No slips, great traction.", "f1fa509f-57a5-4d28-a111-863d470978da"},
		{"030fca7b-f3fe-4638-a7bf-364849daa476", "91183969-4b3a-4692-8ec6-74e2caf131af", "Loving the chunky sole on these Soda Retro 90s. Total nostalgia trip.", "25d45deb-c5ab-449d-b9a8-5f01d061972e"},
		{"8b259397-e77e-473f-a705-544cccc84aa1", "91183969-4b3a-4692-8ec6-74e2caf131af", "My feet breathe so well in the Soda Knit Runners. Perfect for summer.", "54972b8c-b1fd-4cab-8788-53a5fd193a1f"},
		{"f8d66b18-218a-4edc-80f3-0910e0c02996", "91183969-4b3a-4692-8ec6-74e2caf131af", "Soda Pro Skaters holding up well after a week of intense sessions.", "61aad777-8e7c-4135-ab2b-590561bcaea5"},
	}

	log.InfoContext(ctx, "Starting referral-blog seed...")
	for _, b := range blogs {
		_, err := bStore.GetBlog(ctx, b.ID)
		if err == nil {
			log.InfoContext(ctx, "Blog already exists, skipping", "id", b.ID)
			continue
		}

		_, err = bStore.CreateBlog(ctx, db.CreateBlogParams{
			ID:        b.ID,
			AuthorID:  b.AuthorID,
			Content:   b.Content,
			ProductID: b.ProductID,
		})
		if err != nil {
			return fmt.Errorf("creating blog %s: %w", b.ID, err)
		}
		log.InfoContext(ctx, "Created blog", "id", b.ID)
	}
	
	log.InfoContext(ctx, "Blog seed completed successfully")
	return nil
}