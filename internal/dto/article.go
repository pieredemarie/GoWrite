package dto

import "time"

type GetArticleResponse struct {
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	ContentMD string    `json:"content_md"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateArticleRequest struct {
	Title     string `json:"title"`
	ContentMD string `json:"content_md"`
}

type CreateArticleRequest struct {
	Title     string `json:"title"`
	ContentMD string `json:"content_md"`
}

type CreateArticleResponse struct {
	Slug  string `json:"slug"`
	Token string `json:"token"`
}
