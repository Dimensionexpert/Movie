package scanner

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Dimensionexpert/movieLibrary/internal/movie"
)

// ScanMovies walks the directory root and returns all .mp4 files found.
func ScanMovies(root string) ([]movie.MovieFile, error) {
	var found []movie.MovieFile

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".mp4" {
			return nil
		}
		filename := filepath.Base(path)
		title := strings.TrimSuffix(filename, filepath.Ext(path))
		found = append(found, movie.MovieFile{Title: title, Path: path})
		return nil
	})

	return found, err
}
