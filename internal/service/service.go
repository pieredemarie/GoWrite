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

type articleService struct {
	repo ArticleRepository
}

func NewArticleService(repo ArticleRepository) *articleService {
	return &articleService{
		repo: repo,
	}
}

func (s *articleService) CreateArticle(ctx context.Context, title, contentMD string) (slug, token string, error error) {

}

func (s *articleService) GetArticleBySlug(ctx context.Context, slug string) (model.Article, error) {

}

func (s *articleService) UpdateArticle(ctx context.Context, slug, token, title, contentMD string) error {

}

func (s *articleService) DeleteArticle(ctx context.Context, slug, token string) error {

}
