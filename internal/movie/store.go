package movie

import (
	"database/sql"
	"fmt"
	"strconv"

	tmdb "github.com/Dimensionexpert/movieLibrary/internal/TMDB"
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

func (s *Store) GetByID(id int) (Movie, error) {
	var m Movie
	row := s.db.QueryRow(
		"SELECT id, title, file_path, duration_seconds, thumbnail_path, tmdb_id, overview, poster_url, release_year, added_at FROM movies WHERE id = ?",
		id,
	)
	err := row.Scan(&m.ID, &m.Title, &m.FilePath, &m.DurationSeconds, &m.ThumbnailPath, &m.TMDBID, &m.Overview, &m.PosterURL, &m.ReleaseYear, &m.AddedAt)
	if err != nil {
		return Movie{}, fmt.Errorf("error getting movie %d: %w", id, err)
	}
	return m, nil
}

func (s *Store) GetAll() ([]Movie, error) {
	rows, err := s.db.Query("SELECT id, title, file_path, duration_seconds, thumbnail_path, tmdb_id, overview, poster_url, release_year, added_at FROM movies")
	if err != nil {
		return nil, fmt.Errorf("error querying movies: %w", err)
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.FilePath, &m.DurationSeconds, &m.ThumbnailPath, &m.TMDBID, &m.Overview, &m.PosterURL, &m.ReleaseYear, &m.AddedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning movie row: %w", err)
		}
		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating movie rows: %w", err)
	}

	return movies, nil
}

func extractYear(releaseDate string) sql.NullInt64 {
	if len(releaseDate) < 4 {
		return sql.NullInt64{Valid: false}
	}
	year, err := strconv.Atoi(releaseDate[:4])
	if err != nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(year), Valid: true}
}

func (s *Store) UpdateFromTMDB(movieID int, details tmdb.MovieDetails) error {
	_, err := s.db.Exec(
		`UPDATE movies SET overview = ?, poster_url = ?, release_year = ?, duration_seconds = ?, tmdb_id = ? WHERE id = ?`,
		details.Overview, details.PosterPath, extractYear(details.ReleaseDate), details.Runtime*60, details.ID, movieID,
	)
	if err != nil {
		return fmt.Errorf("error updating movie %d from TMDB: %w", movieID, err)
	}

	for _, g := range details.Genres {
		_, err := s.db.Exec(`INSERT INTO genres (id, name) VALUES (?, ?) ON CONFLICT (id) DO NOTHING`, g.ID, g.Name)
		if err != nil {
			return fmt.Errorf("error inserting genre %d: %w", g.ID, err)
		}
		_, err = s.db.Exec(`INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?) ON CONFLICT DO NOTHING`, movieID, g.ID)
		if err != nil {
			return fmt.Errorf("error linking movie %d to genre %d: %w", movieID, g.ID, err)
		}
	}

	return nil
}
