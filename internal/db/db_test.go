package db

import (
	"os"
	"testing"
)

func TestForeignKeysEnforced(t *testing.T) {
	testPath := "test_fk.db"
	defer os.Remove(testPath)

	database, err := OpenDB(testPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer database.Close()

	_, err = database.Exec(
		`INSERT INTO watch_progress (user_id, movie_id, position_seconds) VALUES (?, ?, ?)`,
		999, 1, 0,
	)
	if err == nil {
		t.Fatal("expected foreign key violation, but insert succeeded")
	}
}
