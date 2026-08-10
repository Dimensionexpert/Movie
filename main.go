package main

import (
	"log"
	"net/http"
	"os"

	tmdb "github.com/Dimensionexpert/movieLibrary/internal/TMDB"
	"github.com/Dimensionexpert/movieLibrary/internal/cache"
	"github.com/Dimensionexpert/movieLibrary/internal/db"
	"github.com/Dimensionexpert/movieLibrary/internal/movie"
	"github.com/Dimensionexpert/movieLibrary/internal/scanner"
	"github.com/Dimensionexpert/movieLibrary/web"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		log.Fatal("TMDB_API_KEY not set")
	}
	//temp call
	tmdbClient := tmdb.NewClient(apiKey)
	best, err := tmdbClient.SearchMovie("Inception")
	if err != nil {
		log.Printf("TMDB search failed: %v", err)
	} else {
		log.Printf("best match: [%d] %s (%s)", best.ID, best.Title, best.ReleaseDate)

		details, err := tmdbClient.GetMovieDetails(best.ID)
		if err != nil {
			log.Printf("TMDB details failed: %v", err)
		} else {
			log.Printf("TMDB movie details: %+v", details)

			if err := store.UpdateFromTMDB(1, details); err != nil {
				log.Printf("failed to update movie 1 from TMDB: %v", err)
			} else {
				log.Println("movie 1 successfully enriched from TMDB")
			}
		}
	}

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
