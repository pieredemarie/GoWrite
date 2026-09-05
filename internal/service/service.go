package service

import (
	"context"
	"gowrite/internal/model"
)

type ArticleRepository interface {
	Create(ctx context.Context, article model.Article) error
	Get(ctx context.Context, slug string) (model.Article, error)
	Delete(ctx context.Context, slug string) error
	Update(ctx context.Context, article model.Article) error
}
