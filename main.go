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
		log.Println("TMDB_API_KEY not set, skipping TMDB enrichment")
	}

	err := os.MkdirAll("data", 0755)

	database, err := db.OpenDB("data/movielibrary.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()
	log.Println("Database ready.")

	store := movie.NewStore(database)

	files, err := scanner.ScanMovies("/tmp/testmovies")
	if err != nil {
		log.Printf("scan failed: %v", err)
	} else {
		log.Printf("found %d movie files", len(files))

		removed, err := store.DeleteMissingMovies(files)
		if err != nil {
			log.Printf("failed to reconcile missing movies: %v", err)
		} else {
			log.Printf("removed %d missing movies", removed)
		}

		inserted, err := store.CreateBatch(files)
		if err != nil {
			log.Printf("failed to insert movies: %v", err)
		} else {
			log.Printf("inserted %d movies into database", inserted)
		}
	}

	if apiKey != "" {
		tmdbClient := tmdb.NewClient(apiKey)

		movies, err := store.GetAll()
		if err != nil {
			log.Printf("failed to load movies for enrichment: %v", err)
		} else {
			for _, m := range movies {
				if m.TMDBID.Valid {
					continue
				}

				match, err := tmdbClient.SearchMovie(m.Title)
				if err != nil {
					log.Printf("TMDB search failed for %q: %v", m.Title, err)
					continue
				}

				details, err := tmdbClient.GetMovieDetails(match.ID)
				if err != nil {
					log.Printf("TMDB details failed for %q: %v", m.Title, err)
					continue
				}

				if err := store.UpdateFromTMDB(m.ID, details); err != nil {
					log.Printf("failed to enrich %q: %v", m.Title, err)
					continue
				}

				log.Printf("enriched movie %q from TMDB", m.Title)
			}
		}
	}

	movieCache := cache.NewCache(2)
	mux := web.NewMux(store, movieCache)
	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
