package referralblog

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"soda-interview/business/data/stores/db"
	"soda-interview/foundation/logger"
)

var (
	ErrNotFound = errors.New("blog not found")
)

type Blog struct {
	ID              string
	AuthorID        string
	Content         string
	LinkedProductID string
}

type NewBlog struct {
	AuthorID  string
	Content   string
	ProductID string
}

type Service struct {
	log *logger.Logger
	store Storer
}

func NewService(log *logger.Logger, store Storer) *Service {
	return &Service{
		log: log,
		store: store,
	}
}

func (s *Service) CreateBlog(ctx context.Context, nb NewBlog) (Blog, error) {
	id := uuid.NewString()
	
	dbBlog, err := s.store.CreateBlog(ctx, db.CreateBlogParams{
		ID:        id,
		AuthorID:  nb.AuthorID,
		Content:   nb.Content,
		ProductID: nb.ProductID,
	})
	if err != nil {
		return Blog{}, fmt.Errorf("creating blog: %w", err)
	}

	return toBlog(dbBlog), nil
}

func (s *Service) GetBlog(ctx context.Context, id string) (Blog, error) {
	b, err := s.store.GetBlog(ctx, id)
	if err != nil {
		return Blog{}, fmt.Errorf("querying blog: %w", err)
	}
	return toBlog(b), nil
}

func (s *Service) ListBlogs(ctx context.Context) ([]Blog, error) {
	blogs, err := s.store.ListBlogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing blogs: %w", err)
	}

	out := make([]Blog, len(blogs))
	for i, b := range blogs {
		out[i] = toBlog(b)
	}
	return out, nil
}

func toBlog(dbB db.Blog) Blog {
	return Blog{
		ID:              dbB.ID,
		AuthorID:        dbB.AuthorID,
		Content:         dbB.Content,
		LinkedProductID: dbB.ProductID,
	}
}