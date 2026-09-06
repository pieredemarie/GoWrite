package handler

import (
	"context"
	"gowrite/internal/model"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, title, contentMD string) (articleSlug, token string, error error)
	GetArticleBySlug(ctx context.Context, slug string) (model.Article, error)
	UpdateArticle(ctx context.Context, slug, token, title, contentMD string) error
	DeleteArticle(ctx context.Context, slug, token string) error
}
