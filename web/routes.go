package web

import (
	"embed"
	"net/http"

	"github.com/Dimensionexpert/movieLibrary/internal/cache"
	"github.com/Dimensionexpert/movieLibrary/internal/movie"
)

// staticFiles are embedded so the UI is available regardless of the process
// working directory or which files a production deployment copies alongside the
// binary.
//
//go:embed index.html detail.html style.css app.js detail.js
var staticFiles embed.FS

func NewMux(store *movie.Store, movieCache *cache.Cache) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies/{id}", handleGetMovie(store, movieCache))
	mux.HandleFunc("GET /movies", handleGetMovies(store))
	mux.HandleFunc("GET /movies/{id}/stream", handleStreamMovie(store))

	mux.Handle("GET /", http.FileServer(http.FS(staticFiles)))

	return mux
}
