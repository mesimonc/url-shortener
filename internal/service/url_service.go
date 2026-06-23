package service

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "url-shortener/internal/repository"
)

type URLService struct {
    repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
    return &URLService{repo: repo}
}

func (s *URLService) Shorten(originalURL string) (string, error) {
    code, err := generateCode(6)
    if err != nil {
        return "", fmt.Errorf("generate code: %w", err)
    }

    _, err = s.repo.Save(code, originalURL)
    if err != nil {
        return "", fmt.Errorf("save url: %w", err)
    }

    return code, nil
}

func generateCode(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}