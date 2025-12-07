package referralblog

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
	ErrNotFound = errors.New("blog not found")
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

func (s *Store) CreateBlog(ctx context.Context, params db.CreateBlogParams) (db.Blog, error) {
	b, err := s.q.CreateBlog(ctx, params)
	if err != nil {
		return db.Blog{}, fmt.Errorf("creating blog: %w", err)
	}
	return b, nil
}

func (s *Store) GetBlog(ctx context.Context, id string) (db.Blog, error) {
	b, err := s.q.GetBlog(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return db.Blog{}, ErrNotFound
		}
		return db.Blog{}, fmt.Errorf("querying blog: %w", err)
	}
	return b, nil
}

func (s *Store) ListBlogs(ctx context.Context) ([]db.Blog, error) {
	blogs, err := s.q.ListBlogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing blogs: %w", err)
	}
	return blogs, nil
}