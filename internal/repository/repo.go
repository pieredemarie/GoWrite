package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gowrite/internal/model"
	"time"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(dsn string) (*PostgresRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &PostgresRepo{
		db: db,
	}, nil
}

func (p *PostgresRepo) Create(ctx context.Context, article model.Article) error {
	query := `
		INSERT INTO articles (id,slug, title, contentmd, contenthtml, 		  authortokenhash, created_at, updated_at)
		VALUES ($1, $2,$3,$4,$5, $6,$7) `

	_, err := p.db.ExecContext(
		ctx,
		query,
		article.ID,
		article.Slug,
		article.Title,
		article.ContentMD,
		article.ContentHTML,
		article.AuthorTokenHash,
		article.CreatedAt,
		article.UpdatedAt,
	)
	return err
}

func (p *PostgresRepo) Delete(ctx context.Context, slug string) error {
	query := `
		DELETE FROM articles WHERE slug = $1
	`

	_, err := p.db.ExecContext(
		ctx,
		query,
		slug,
	)

	return err
}

func (p *PostgresRepo) Get(ctx context.Context, slug string) (model.Article, error) {
	query := `
		SELECT id, slug, title, contentmd, contenthtml, authortokenhash, created_at, updated_at
		WHERE id = $1
	`

	var article model.Article

	err := p.db.QueryRowContext(ctx, query).Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.ContentMD,
		&article.ContentHTML,
		&article.AuthorTokenHash,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Article{}, fmt.Errorf("article with slug %s not found %w", slug, err)
		}
		return model.Article{}, err
	}

	return article, nil
}

func (p *PostgresRepo) Update(ctx context.Context) error {
	// TODO: decide how to update article
	// using the whole struct (about 24 bytes) or just contentHTML (or md?)

	return nil
}

func (p *PostgresRepo) Close() {
	p.db.Close()
}
