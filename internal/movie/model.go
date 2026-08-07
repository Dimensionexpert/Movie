package movie

import (
	"database/sql"
	"time"
)

type Movie struct {
	ID              int
	Title           string
	FilePath        string
	DurationSeconds sql.NullInt64
	ThumbnailPath   sql.NullString
	TMDBID          sql.NullInt64
	Overview        sql.NullString
	PosterURL       sql.NullString
	ReleaseYear     sql.NullInt64
	AddedAt         time.Time
}

// MovieFile represents a discovered movie file on disk, found by the scanner
// before it's been inserted into the database.
type MovieFile struct {
	Title string
	Path  string
}
