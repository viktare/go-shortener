package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/viktare/go-shortener/model"
)

type UrlRepository struct {
	db *sql.DB
}

func NewUrlRepository(db *sql.DB) *UrlRepository {
	return &UrlRepository{db: db}
}

func (r *UrlRepository) Create(ctx context.Context, url model.Url) (model.Url, error) {
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO urls (short_url, original_url) 
			VALUES ($1, $2) 
			ON CONFLICT (short_url) DO NOTHING
			RETURNING short_url, original_url, created_at`,
		url.ShortUrl,
		url.OriginalUrl,
	).Scan(&url.ShortUrl, &url.OriginalUrl, &url.CreatedAt)

	if err == sql.ErrNoRows {
		// conflict happened, just fetch the existing record
		return r.FindByShortUrl(ctx, url.ShortUrl)
	}
	if err != nil {
		return model.Url{}, fmt.Errorf("failed to insert url: %w", err)
	}

	return url, nil
}

func (r *UrlRepository) FindByShortUrl(ctx context.Context, shortUrl string) (model.Url, error) {
	var url model.Url

	err := r.db.QueryRowContext(
		ctx,
		`SELECT short_url, original_url, created_at FROM urls WHERE short_url = $1`,
		shortUrl,
	).Scan(&url.ShortUrl, &url.OriginalUrl, &url.CreatedAt)

	if err == sql.ErrNoRows {
		return model.Url{}, fmt.Errorf("url not found")
	}
	if err != nil {
		return model.Url{}, fmt.Errorf("failed to find url: %w", err)
	}

	return url, nil
}
