package main

import (
	"moviepin/middleware"
	"moviepin/routes"
	"net/http"
)

func main() {
	mux := routes.NewServeMux()

	loggedMux := middleware.Logger(mux)

	http.ListenAndServe(":4545", loggedMux)
}
