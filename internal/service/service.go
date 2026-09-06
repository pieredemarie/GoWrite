package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"gowrite/internal/customerr"
	"gowrite/internal/model"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
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
func RandomID(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random id: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func RenderMarkdown(md string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}
	policy := bluemonday.UGCPolicy()
	safeHTML := policy.Sanitize(buf.String())
	return safeHTML, nil
}

func (s *articleService) CreateArticle(ctx context.Context, title, contentMD string) (articleSlug, token string, error error) {
	if title == "" {
		return "", "", customerr.ErrTitleEmpty
	}
	if len(strings.Fields(contentMD)) > 1000 {
		return "", "", customerr.ErrMaxLimitExceeded
	}
	suffix, err := RandomID(8)
	if err != nil {
		return "", "", err
	}

	articleSlug = slug.Make(title) + "-" + suffix

	token, err = GenerateToken()
	if err != nil {
		return "", "", err
	}

	tokenHash := HashToken(token)

	html, err := RenderMarkdown(contentMD)
	if err != nil {
		return "", "", err
	}

	now := time.Now()
	article := model.Article{
		Slug:            articleSlug,
		Title:           title,
		ContentMD:       contentMD,
		ContentHTML:     html,
		AuthorTokenHash: tokenHash,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Create(ctx, article); err != nil {
		return "", "", err
	}

	return articleSlug, token, nil
}

func (s *articleService) GetArticleBySlug(ctx context.Context, slug string) (model.Article, error) {

}

func (s *articleService) UpdateArticle(ctx context.Context, slug, token, title, contentMD string) error {

}

func (s *articleService) DeleteArticle(ctx context.Context, slug, token string) error {

}
