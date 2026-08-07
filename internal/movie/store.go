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
//
// TODO: this will be the final step of the upload pipeline once built:
// HTTP handler (receive multipart upload) -> validate/sanitize (check
// extension + verify actual codec/container via ffprobe, don't trust
// the extension alone) -> move file to permanent media dir -> Create.
// Scanning (CreateBatch) and uploading are two separate producers that
// both funnel into this same insert path.
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

func (s *Store) CreateBatch(files []MovieFile) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction : %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO movies (title, file_path) VALUES (?, ?) ON CONFLICT (file_path) DO NOTHING")
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for _, f := range files {
		result, err := stmt.Exec(f.Title, f.Path)
		if err != nil {
			return inserted, fmt.Errorf("error inserting movie %q: %w", f.Title, err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return inserted, fmt.Errorf("error checking rows affected for %q: %w", f.Title, err)
		}
		if rows > 0 {
			inserted++
		}
	}

	if err := tx.Commit(); err != nil {
		return inserted, fmt.Errorf("error Comitting transcation: %w", err)
	}

	return inserted, nil
}
