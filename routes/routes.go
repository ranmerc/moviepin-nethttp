package routes

import (
	"moviepin/handlers"
	"net/http"
)

// Returns a mux with all routes added.
func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/movies", handlers.NewMoviesHandler())
	mux.Handle("/movies/", handlers.NewMoviesHandler())

	return mux
}
