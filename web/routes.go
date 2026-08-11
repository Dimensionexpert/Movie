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
	mux.HandleFunc("GET /movies/{id}/stream", handleStreamMovie(store))

	mux.Handle("GET /", http.FileServer(http.Dir("web")))

	return mux
}
