package repository

// import (
//     "database/sql"
//     "fmt"
// 	"time"

//     _ "github.com/lib/pq"
// )

import (
	"time"

	"gorm.io/gorm"
)

// URL represents a shortened URL record in the database.
type URL struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"uniqueIndex;not null;size:10"`
	OriginalURL string    `gorm:"not null"`
	Clicks      int64     `gorm:"default:0"`
	CreatedAt   time.Time
	ExpiresAt   *time.Time `gorm:"index"`
}

type URLRepository struct {
    //db *sql.DB
	db *gorm.DB
}

// func NewURLRepository(db *sql.DB) *URLRepository {
//     return &URLRepository{db: db}
// }

// NewURLRepository creates a new URLRepository with the given GORM database.
func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{db: db}
}

// func (r *URLRepository) Save(code, originalURL string) (*URL, error) {
//     query := `
//         INSERT INTO urls (code, original_url)
//         VALUES ($1, $2)
//         RETURNING id, code, original_url, clicks, created_at
//     `
//     url := &URL{}
//     err := r.db.QueryRow(query, code, originalURL).Scan(
//         &url.ID, &url.Code, &url.OriginalURL, &url.Clicks, &url.CreatedAt,
//     )
//     if err != nil {
//         return nil, fmt.Errorf("save url: %w", err)
//     }
//     return url, nil
// }

// Save inserts a new URL record into the database.
func (r *URLRepository) Save(code, originalURL string, expiresAt *time.Time) (*URL, error) {
	url := &URL{Code: code, OriginalURL: originalURL, ExpiresAt: expiresAt}
	result := r.db.Create(url)
	return url, result.Error
}

// func (r *URLRepository) FindByCode(code string) (*URL, error) {
//     query := `SELECT id, code, original_url, clicks, created_at FROM urls WHERE code = $1`
//     url := &URL{}
//     err := r.db.QueryRow(query, code).Scan(
//         &url.ID, &url.Code, &url.OriginalURL, &url.Clicks, &url.CreatedAt,
//     )
//     if err == sql.ErrNoRows {
//         return nil, nil
//     }
//     if err != nil {
//         return nil, err
//     }
//     return url, nil
// }

// FindByCode retrieves a URL record by its short code.
func (r *URLRepository) FindByCode(code string) (*URL, error) {
	var url URL
	result := r.db.Where("code = ?", code).First(&url)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &url, result.Error
}

// IncrementClicks increments the click count for the given short code.
// func (r *URLRepository) IncrementClicks(code string) error {
//     _, err := r.db.Exec(`UPDATE urls SET clicks = clicks + 1 WHERE code = $1`, code)
//     return err
// }

// IncrementClicks increments the click count for the given short code.
func (r *URLRepository) IncrementClicks(code string) error {
	return r.db.Model(&URL{}).Where("code = ?", code).
		UpdateColumn("clicks", gorm.Expr("clicks + 1")).Error
}