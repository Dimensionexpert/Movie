package movie

import (
	"os"
	"testing"

	"github.com/Dimensionexpert/movieLibrary/internal/db"
)

func TestCreateBatchSkipsDuplicates(t *testing.T) {
	testPath := "test_CB.db"
	defer os.Remove(testPath)
	database, err := db.OpenDB(testPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	store := NewStore(database)
	_, err = store.Create("Inception.mp4", "/tmp/testmovies/Inception.mp4")
	if err != nil {
		t.Fatalf("failed to insert seed movie: %v", err)
	}

	files := []MovieFile{
		{Title: "New Movie One", Path: "/tmp/testmovies/NewOne.mp4"},
		{Title: "New Movie Two", Path: "/tmp/testmovies/NewTwo.mp4"},
		{Title: "Inception.mp4", Path: "/tmp/testmovies/Inception.mp4"}, // duplicate of seed
	}

	inserted, err := store.CreateBatch(files)
	if err != nil {
		t.Fatalf("expected CreateBatch to succeed despite duplicate, got error: %v", err)
	}

	if inserted != 2 {
		t.Fatalf("expected 2 new rows inserted (duplicate skipped), got %d", inserted)
	}

	var count int
	row := database.QueryRow("SELECT COUNT(*) FROM movies")
	if err := row.Scan(&count); err != nil {
		t.Fatalf("failed to count movies: %v", err)
	}

	// seed (1) + 2 new = 3 total; duplicate should not have added a 4th row
	if count != 3 {
		t.Fatalf("expected 3 total rows in database, got %d", count)
	}
}
