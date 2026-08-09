package movie

import (
	"database/sql"
	"encoding/json"
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

func (m Movie) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"id":       m.ID,
		"title":    m.Title,
		"filePath": m.FilePath,
		"addedAt":  m.AddedAt,
	}

	if m.DurationSeconds.Valid {
		data["durationSeconds"] = m.DurationSeconds.Int64
	} else {
		data["durationSeconds"] = nil
	}

	if m.ThumbnailPath.Valid {
		data["thumbnailPath"] = m.ThumbnailPath.String
	} else {
		data["thumbnailPath"] = nil
	}

	if m.TMDBID.Valid {
		data["tmdbId"] = m.TMDBID.Int64
	} else {
		data["tmdbId"] = nil
	}

	if m.Overview.Valid {
		data["overview"] = m.Overview.String
	} else {
		data["overview"] = nil
	}

	if m.PosterURL.Valid {
		data["posterUrl"] = m.PosterURL.String
	} else {
		data["posterUrl"] = nil
	}

	if m.ReleaseYear.Valid {
		data["releaseYear"] = m.ReleaseYear.Int64
	} else {
		data["releaseYear"] = nil
	}

	return json.Marshal(data)
}

// MovieFile represents a discovered movie file on disk, found by the scanner
// before it's been inserted into the database.
type MovieFile struct {
	Title string
	Path  string
}
