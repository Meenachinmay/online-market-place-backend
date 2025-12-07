package referralblog

import (
	"context"

	"soda-interview/business/data/stores/db"
)

type Storer interface {
	CreateBlog(ctx context.Context, params db.CreateBlogParams) (db.Blog, error)
	GetBlog(ctx context.Context, id string) (db.Blog, error)
	ListBlogs(ctx context.Context) ([]db.Blog, error)
}
