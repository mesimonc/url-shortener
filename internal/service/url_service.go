package service

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "url-shortener/internal/repository"
)

type URLService struct {
    repo  *repository.URLRepository
    cache *repository.Cache
}

// NewURLService creates a new URLService with the given repository and cache.
func NewURLService(repo *repository.URLRepository, cache *repository.Cache) *URLService {
    return &URLService{repo: repo, cache: cache}
}

// Shorten generates a short code for the given URL and persists it.
func (s *URLService) Shorten(originalURL string) (string, error) {
    code, err := generateCode(6)
    if err != nil {
        return "", fmt.Errorf("generate code: %w", err)
    }

    _, err = s.repo.Save(code, originalURL)
    if err != nil {
        return "", fmt.Errorf("save url: %w", err)
    }

    _ = s.cache.Set(context.Background(), code, originalURL)

    return code, nil
}

// Resolve looks up the original URL for the given short code.
// It checks the cache first, falling back to the database on a miss.
func (s *URLService) Resolve(code string) (string, error) {
    // Check cache first
    if cached, err :=