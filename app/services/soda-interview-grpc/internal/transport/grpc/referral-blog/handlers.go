package referralblog

import (
	"context"

	"soda-interview/business/core/referral-blog"
	referralblogv1 "soda-interview/foundation/proto/referral-blog/v1"
)

type Handler struct {
	referralblogv1.UnimplementedBlogServiceServer
	Service *referralblog.Service
}

func (h *Handler) CreateBlog(ctx context.Context, req *referralblogv1.CreateBlogRequest) (*referralblogv1.Blog, error) {
	nb := referralblog.NewBlog{
		AuthorID:  req.AuthorId,
		Content:   req.Content,
		ProductID: req.ProductId,
	}

	b, err := h.Service.CreateBlog(ctx, nb)
	if err != nil {
		return nil, err
	}

	return &referralblogv1.Blog{
		Id:              b.ID,
		AuthorId:        b.AuthorID,
		Content:         b.Content,
		LinkedProductId: b.LinkedProductID,
	}, nil
}

func (h *Handler) GetBlog(ctx context.Context, req *referralblogv1.BlogRequest) (*referralblogv1.Blog, error) {
	b, err := h.Service.GetBlog(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &referralblogv1.Blog{
		Id:              b.ID,
		AuthorId:        b.AuthorID,
		Content:         b.Content,
		LinkedProductId: b.LinkedProductID,
	}, nil
}

func (h *Handler) ListBlogs(ctx context.Context, _ *referralblogv1.Empty) (*referralblogv1.BlogList, error) {
	blogs, err := h.Service.ListBlogs(ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*referralblogv1.Blog, len(blogs))
	for i, b := range blogs {
		list[i] = &referralblogv1.Blog{
			Id:              b.ID,
			AuthorId:        b.AuthorID,
			Content:         b.Content,
			LinkedProductId: b.LinkedProductID,
		}
	}

	return &referralblogv1.BlogList{Blogs: list}, nil
}