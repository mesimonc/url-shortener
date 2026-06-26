package repository

// import (
//     "database/sql"
//     "fmt"

//     _ "github.com/lib/pq"
// )

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
// func NewPostgres(databaseURL string) (*sql.DB, error) {
//     db, err := sql.Open("postgres", databaseURL)
//     if err != nil {
//         return nil, fmt.Errorf("open postgres: %w", err)
//     }

//     if err := db.Ping(); err != nil {
//         return nil, fmt.Errorf("ping postgres: %w", err)
//     }

//     return db, nil
// }

// NewPostgres creates a new GORM database connection and runs auto migration.
func NewPostgres(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// func Migrate(db *sql.DB) error {
//     _, err := db.Exec(`
//         CREATE TABLE IF NOT EXISTS urls (
//             id          BIGSERIAL PRIMARY KEY,
//             code        VARCHAR(10) UNIQUE NOT NULL,
//             original_url TEXT NOT NULL,
//             clicks      BIGINT DEFAULT 0,
//             created_at  TIMESTAMPTZ DEFAULT NOW()
//         );
//     `)
//     return err
// }

// Migrate runs auto migration for all models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&URL{})
}
