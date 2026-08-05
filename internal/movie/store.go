package movie

import (
	"database/sql"
	"fmt"
)

// Store handles database operations for movies.
type Store struct {
	db *sql.DB
}

// NewStore creates a new movie Store backed by the given database connection.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Create inserts a new movie with the given title and file path,
// returning the auto-generated ID of the inserted row.
func (s *Store) Create(title, filePath string) (int64, error) {
	result, err := s.db.Exec(
		"INSERT INTO movies (title, file_path) VALUES (?, ?)",
		title, filePath,
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting movie: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting inserted movie id: %w", err)
	}

	return id, nil
}
