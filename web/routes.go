package web

import (
	"net/http"

	"github.com/Dimensionexpert/movieLibrary/internal/cache"
	"github.com/Dimensionexpert/movieLibrary/internal/movie"
)

func NewMux(store *movie.Store, movieCache *cache.Cache) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies/{id}", handleGetMovie(store, movieCache))
	mux.HandleFunc("GET /movies", handleGetMovies(store))
	return mux
}
