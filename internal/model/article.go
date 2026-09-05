package model

import "time"

type Article struct {
	ID              int       `json:"-"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	ContentMD       string    `json:"content_md,omitempty"`
	ContentHTML     string    `json:"content_html"`
	AuthorTokenHash string    `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
