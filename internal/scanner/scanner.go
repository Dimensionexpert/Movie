package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// ScanMovies walks the directory root and prints the title and path of .mp4 files.
func ScanMovies(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Check for .mp4 extension
		if filepath.Ext(path) != ".mp4" {
			return nil
		}

		// Get filename and remove the extension
		filename := filepath.Base(path)
		title := strings.TrimSuffix(filename, filepath.Ext(path))

		fmt.Println(title, path)
		return nil
	})
}
