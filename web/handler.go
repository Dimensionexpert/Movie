package web

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Dimensionexpert/movieLibrary/internal/cache"
	"github.com/Dimensionexpert/movieLibrary/internal/movie"
)

func handleGetMovie(store *movie.Store, movieCache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqIdStr := r.PathValue("id")
		reqIdInt, err := strconv.Atoi(reqIdStr)
		if err != nil {
			http.Error(w, "Invalid movie ID format", http.StatusBadRequest)
			return
		}

		m, ok := movieCache.Get(reqIdInt)
		if ok {
			log.Printf("cache hit for movie %d", reqIdInt)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(m)
			return
		}

		log.Printf("cache miss for movie %d, querying database", reqIdInt)
		md, err := store.GetByID(reqIdInt)
		if err != nil {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}

		movieCache.Put(reqIdInt, md)
		log.Printf("database hit for movie %d, cached for future requests", reqIdInt)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(md)

	}
}
func handleGetMovies(store *movie.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies, err := store.GetAll()
		if err != nil {
			http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movies)
	}
}

func handleStreamMovie(store *movie.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqIdStr := r.PathValue("id")
		reqIdInt, err := strconv.Atoi(reqIdStr)
		if err != nil {
			http.Error(w, "Invalid movie ID format", http.StatusBadRequest)
			return
		}

		m, err := store.GetByID(reqIdInt)
		if err != nil {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}

		file, err := os.Open(m.FilePath)
		if err != nil {
			http.Error(w, "Could not open video file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			http.Error(w, "Could not stat video file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "video/mp4")
		http.ServeContent(w, r, m.Title, stat.ModTime(), file)
	}
}
