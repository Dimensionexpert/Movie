package main

import (
	"log"
	"net/http"

	"github.com/Dimensionexpert/movieLibrary/internal/cache"
	"github.com/Dimensionexpert/movieLibrary/internal/db"
	"github.com/Dimensionexpert/movieLibrary/internal/movie"
	"github.com/Dimensionexpert/movieLibrary/internal/scanner"
	"github.com/Dimensionexpert/movieLibrary/web"
)

func main() {
	database, err := db.OpenDB("data/movielibrary.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()
	log.Println("Database ready.")

	store := movie.NewStore(database)

	files, err := scanner.ScanMovies("/tmp/testmovies")
	if err != nil {
		log.Fatalf("scan failed: %v", err)
	}
	log.Printf("found %d movie files", len(files))

	inserted, err := store.CreateBatch(files)
	if err != nil {
		log.Fatalf("failed to insert movies: %v", err)
	}
	log.Printf("inserted %d movies into database", inserted)

	movieCache := cache.NewCache(2)

	mux := web.NewMux(store, movieCache)
	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
