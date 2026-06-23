package repository

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

type URL struct {
    ID          int64
    Code        string
    OriginalURL string
    Clicks      int64
}

type URLRepository struct {
    db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
    return &URLRepository{db: db}
}

func (r *URLRepository) Save(code, originalURL string) (*URL, error) {
    query := `
        INSERT INTO urls (code, original_url)
        VALUES ($1, $2)
        RETURNING id, code, original_url, clicks
    `
    url := &URL{}
    err := r.db.QueryRow(query, code, originalURL).Scan(
        &url.ID, &url.Code, &url.OriginalURL, &url.Clicks,
    )
    if err != nil {
        return nil, fmt.Errorf("save url: %w", err)
    }
    return url, nil
}

func (r *URLRepository) FindByCode(code string) (*URL, error) {
    query := `SELECT id, code, original_url, clicks FROM urls WHERE code = $1`
    url := &URL{}
    err := r.db.QueryRow(query, code).Scan(
        &url.ID, &url.Code, &url.OriginalURL, &url.Clicks,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return url, nil
}