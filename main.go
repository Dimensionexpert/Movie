package main

import (
	"log"

	"github.com/Dimensionexpert/movieLibrary/internal/db"
	"github.com/Dimensionexpert/movieLibrary/internal/scanner"
)

func main() {
	database, err := db.OpenDB("data/movielibrary.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	log.Println("Database ready.")

	if err := scanner.ScanMovies("/tmp/testmovies"); err != nil {
		log.Fatalf("scan failed: %v", err)
	}
}
