package web

import (
	"encoding/json"
	"log"
	"net/http"
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
