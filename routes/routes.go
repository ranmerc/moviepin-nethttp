package routes

import (
	"moviepin/db"
	"moviepin/db/movies"
	"moviepin/handlers"
	"net/http"
)

// Returns a mux with all routes added.
func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	moviesDB := movies.NewMovie(db.DB)

	mux.Handle("/movies", handlers.NewMoviesHandler(moviesDB))
	mux.Handle("/movies/", handlers.NewMoviesHandler(moviesDB))

	return mux
}
