package repository

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

func NewPostgres(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, fmt.Errorf("open postgres: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("ping postgres: %w", err)
    }

    return db, nil
}

func Migrate(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS urls (
            id          BIGSERIAL PRIMARY KEY,
            code        VARCHAR(10) UNIQUE NOT NULL,
            original_url TEXT NOT NULL,
            clicks      BIGINT DEFAULT 0,
            created_at  TIMESTAMPTZ DEFAULT NOW()
        );
    `)
    return err
}
