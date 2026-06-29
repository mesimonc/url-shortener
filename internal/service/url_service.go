package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"url-shortener/internal/repository"
	"time"
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
// If customCode is provided, it will be used instead of a random code.
// If expiresInDays is 0, the URL never expires.
func (s *URLService) Shorten(originalURL, customCode string, expiresInDays int) (string, error) {
	code := customCode
	if code == "" {
		var err error
		code, err = generateCode(6)
		if err != nil {
			return "", fmt.Errorf("generate code: %w", err)
		}
	}

	var expiresAt *time.Time
	if expiresInDays > 0 {
		t := time.Now().AddDate(0, 0, expiresInDays)
		expiresAt = &t
	}

	_, err := s.repo.Save(code, originalURL, expiresAt)
	if err != nil {
    if strings.Contains(err.Error(), "duplicate") {
        return "", fmt.Errorf("code already taken")
    }
    return "", fmt.Errorf("save url: %w", err)
}

	_ = s.cache.Set(context.Background(), code, originalURL)

	return code, nil
}

// Resolve looks up the original URL for the given short code.
// It checks the cache first, falling back to the database on a miss.
func (s *URLService) Resolve(code string) (string, error) {
	// Check cache first
	if cached, err := s.cache.Get(context.Background(), code); err == nil && cached != "" {
		go s.repo.IncrementClicks(code)
		return cached, nil
	}

	// Cache miss, query database
	url, err := s.repo.FindByCode(code)
	if err != nil {
		return "", fmt.Errorf("find by code: %w", err)
	}
	if url == nil {
		return "", nil
	}

	// Check if URL has expired
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
    	return "", nil
	}

	// Backfill cache
	_ = s.cache.Set(context.Background(), code, url.OriginalURL)

	// Increment click count asynchronously
	go s.repo.IncrementClicks(code)

	return url.OriginalURL, nil
}

// GetStats returns the URL record for the given short code.
func (s *URLService) GetStats(code string) (*repository.URL, error) {
	return s.repo.FindByCode(code)
}

// generateCode generates a cryptographically random URL-safe string of the given length.
func generateCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}