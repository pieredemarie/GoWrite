package handler

import (
	"context"
	"encoding/json"
	"errors"
	"gowrite/internal/customerr"
	"gowrite/internal/dto"
	"gowrite/internal/model"
	"net/http"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, title, contentMD string) (articleSlug, token string, error error)
	GetArticleBySlug(ctx context.Context, slug string) (model.Article, error)
	UpdateArticle(ctx context.Context, slug, token, title, contentMD string) error
	DeleteArticle(ctx context.Context, slug, token string) error
}

type ArticleHandler struct {
	service ArticleService
}

func NewArticleHandler(service ArticleService) *ArticleHandler {
	return &ArticleHandler{
		service: service,
	}
}

func (a *ArticleHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateArticleRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	slug, token, err := a.service.CreateArticle(ctx, req.Title, req.ContentMD)
	switch {
	case err == nil:
	case errors.Is(err, customerr.ErrTitleEmpty), errors.Is(err, customerr.ErrMaxLimitExceeded):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	resp := dto.CreateArticleResponse{Slug: slug, Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (a *ArticleHandler) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	ctx := r.Context()

	article, err := a.service.GetArticleBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, customerr.ErrNotFoundOrForbidden) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	resp := dto.GetArticleResponse{
		Slug:      article.Slug,
		Title:     article.Title,
		ContentMD: article.ContentMD,
		CreatedAt: article.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (a *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	ctx := r.Context()

	var req dto.UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}
	token := r.Header.Get("X-Edit-Token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	err := a.service.UpdateArticle(ctx, slug, token, req.Title, req.ContentMD)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, customerr.ErrTitleEmpty), errors.Is(err, customerr.ErrMaxLimitExceeded):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, customerr.ErrNotFoundOrForbidden):
		http.Error(w, "not found  or forbidden", http.StatusForbidden)
	default:
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func (a *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	ctx := r.Context()

	token := r.Header.Get("X-Edit-Token")
	if token == "" {
		http.Error(w, "missing edit token", http.StatusUnauthorized)
		return
	}

	err := a.service.DeleteArticle(ctx, slug, token)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, customerr.ErrNotFoundOrForbidden):
		http.Error(w, "not found or forbidden", http.StatusForbidden)
	default:
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}
