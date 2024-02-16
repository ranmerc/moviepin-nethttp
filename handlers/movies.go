package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"moviepin/db"
	"moviepin/models"
	"moviepin/utils"
)

const (
	// ErrRead is returned when failed to read result.
	ErrRead = "failed to read result"

	// ErrNotExists is returned when movie does not exist.
	ErrNotExists = "movie does not exist"

	// ErrFailedToGetMovie is returned when failed to get movie.
	ErrFailedToGetMovie = "failed to get movie"

	// ErrFailedToGetMovies is returned when failed to get movies.
	ErrFailedToGetMovies = "failed to get movies"

	// ErrFailedToAddMovie is returned when failed to add movie.
	ErrFailedToAddMovie = "failed to add movie"

	// ErrFailedToUpdateMovie is returned when failed to update movie.
	ErrFailedToUpdateMovie = "failed to update movie"

	// ErrFailedToDeleteMovie is returned when failed to delete movie.
	ErrFailedToDeleteMovie = "failed to delete movie"

	// ErrFailedToReplaceMovies is returned when failed to replace movies.
	ErrFailedToReplaceMovies = "failed to replace movies"
)

type MoviesHandler struct{}

// Responds with all the movies.
func (mh *MoviesHandler) getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := db.GetMovies()

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovies, http.StatusInternalServerError)
		return
	}

	moviesJson, err := json.Marshal(movies)

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovies, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(moviesJson)
}

// Responds with details of particular movie.
func (mh *MoviesHandler) getMovie(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("rating") == "true" {
		mh.getMovieRating(w, r)
		return
	}

	id, err := utils.GetIDFromPath(r.URL.Path)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = utils.Validate.Var(id, "required,uuid")

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovie, http.StatusBadRequest)
		return
	}

	movie, err := db.GetMovie(id)

	if err == db.ErrNotExists {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovie, http.StatusInternalServerError)
		return
	}

	movieJson, err := json.Marshal(movie)

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovie, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(movieJson)
}

// Adds list of movies sent in request.
func (mh *MoviesHandler) postMovies(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.Logger.Print(err)
		http.Error(w, ErrFailedToAddMovie, http.StatusInternalServerError)
		return
	}

	var movies []models.Movie

	if err = json.Unmarshal(body, &movies); err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToAddMovie, http.StatusBadRequest)
		return
	}

	for _, movie := range movies {
		if err = utils.Validate.Struct(movie); err != nil {
			utils.Logger.Println(err)
			http.Error(w, ErrFailedToAddMovie, http.StatusBadRequest)
			return
		}
	}

	type MovieStatus struct {
		Status string
		Movie  models.Movie
	}

	status := make(chan MovieStatus)

	for _, movie := range movies {
		go func(movie models.Movie) {
			if err := db.AddMovie(movie); err != nil {
				utils.Logger.Print(err)

				status <- MovieStatus{
					Status: "failed",
					Movie:  movie,
				}
				return
			}

			status <- MovieStatus{
				Status: "success",
				Movie:  movie,
			}
		}(movie)
	}

	type postMoviesResponse struct {
		AddedMovies  []*models.Movie `json:"added_movies,omitempty"`
		FailedMovies []*models.Movie `json:"failed_movies,omitempty"`
	}

	response := &postMoviesResponse{}

	for range movies {
		result := <-status

		switch result.Status {
		case "failed":
			response.FailedMovies = append(response.FailedMovies, &result.Movie)
		case "success":
			response.AddedMovies = append(response.AddedMovies, &result.Movie)
		}
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToAddMovie, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

// Updates a particular movie.
func (mh *MoviesHandler) patchMovies(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r.URL.Path)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = utils.Validate.Var(id, "required,uuid")

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.Logger.Print(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusInternalServerError)
		return
	}

	var partialMovie map[string]interface{}

	if err = json.Unmarshal(body, &partialMovie); err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
		return
	}

	existingMovie, err := db.GetMovie(id)

	if err != nil {
		// If movie does not exist.
		if err == db.ErrNotExists {
			http.NotFound(w, r)
			return
		}

		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusInternalServerError)
		return
	}

	// Update fields of existing movie with request.
	for key, value := range partialMovie {
		switch key {
		case "title":
			if title, ok := value.(string); ok {
				existingMovie.Title = title
			} else {
				utils.Logger.Printf("failed to assert type for field title")
				http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
				return
			}
		case "release_date":
			if releaseDate, ok := value.(string); ok {
				if parsedReleaseDate, err := time.Parse(time.RFC3339, releaseDate); err == nil {
					existingMovie.ReleaseDate = parsedReleaseDate
				} else {
					utils.Logger.Printf("failed to parse release_date: %v", err)
					http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
					return
				}
			} else {
				utils.Logger.Printf("failed to assert type for field release_date")
				http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
				return
			}
		case "genre":
			if genre, ok := value.(string); ok {
				existingMovie.Genre = genre
			} else {
				utils.Logger.Printf("failed to assert type for field genre")
				http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
				return
			}
		case "director":
			if director, ok := value.(string); ok {
				existingMovie.Director = director
			} else {
				utils.Logger.Printf("failed to assert type for field director")
				http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
				return
			}
		case "description":
			if description, ok := value.(string); ok {
				existingMovie.Description = description
			} else {
				utils.Logger.Printf("failed to assert type for field description")
				http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
				return
			}
		}
	}

	if err := utils.Validate.Struct(existingMovie); err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
		return
	}

	err = db.UpdateMovie(id, *existingMovie)

	if err != nil {
		if err == db.ErrNotExists {
			http.NotFound(w, r)
			return
		}

		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Deletes a particular movie.
func (mh *MoviesHandler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r.URL.Path)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = utils.Validate.Var(id, "required,uuid")

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToDeleteMovie, http.StatusBadRequest)
		return
	}

	err = db.DeleteMovie(id)

	if err != nil {
		if err == db.ErrNotExists {
			http.NotFound(w, r)
			return
		}

		utils.Logger.Println(err)
		http.Error(w, ErrFailedToDeleteMovie, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Updates a particular movie.
func (mh *MoviesHandler) putMovie(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r.URL.Path)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = utils.Validate.Var(id, "required,uuid")

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.Logger.Print(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusInternalServerError)
	}

	var movie models.Movie

	if err = json.Unmarshal(body, &movie); err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
		return
	}

	if err = utils.Validate.Struct(movie); err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusBadRequest)
		return
	}

	err = db.UpdateMovie(id, movie)

	if err != nil {
		if err == db.ErrNotExists {
			http.NotFound(w, r)
			return
		}

		utils.Logger.Println(err)
		http.Error(w, ErrFailedToUpdateMovie, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Updates whole collection of movies.
func (mh *MoviesHandler) putMovies(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.Logger.Print(err)
		http.Error(w, ErrFailedToReplaceMovies, http.StatusInternalServerError)
		return
	}

	var movies []*models.Movie

	if err := json.Unmarshal(body, &movies); err != nil {
		utils.Logger.Print(err)
		http.Error(w, ErrFailedToReplaceMovies, http.StatusBadRequest)
		return
	}

	for _, movie := range movies {
		if err := utils.Validate.Struct(movie); err != nil {
			utils.Logger.Print(err)
			http.Error(w, ErrFailedToReplaceMovies, http.StatusBadRequest)
			return
		}
	}

	err = db.ReplaceMovies(movies)

	if err != nil {
		utils.Logger.Print(err)
		http.Error(w, ErrFailedToReplaceMovies, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Responds with movie details along with its rating.
func (mh *MoviesHandler) getMovieRating(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromPath(r.URL.Path)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = utils.Validate.Var(id, "required,uuid")

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovie, http.StatusBadRequest)
		return
	}

	_, err = db.GetMovie(id)

	if err == db.ErrNotExists {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovie, http.StatusInternalServerError)
		return
	}

	review, err := db.GetMovieRating(id)

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovie, http.StatusInternalServerError)
		return
	}

	err = utils.Validate.Struct(review)

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovie, http.StatusInternalServerError)
		return
	}

	reviewJson, err := json.Marshal(review)

	if err != nil {
		utils.Logger.Println(err)
		http.Error(w, ErrFailedToGetMovie, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(reviewJson)
}

// Responds with allowed methods.
func (mh *MoviesHandler) Options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
	w.WriteHeader(http.StatusNoContent)
}

func (mh *MoviesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isCollectionPath := r.URL.Path == "/movies" || r.URL.Path == "/movies/"

	switch r.Method {
	case http.MethodGet:
		if isCollectionPath {
			mh.getMovies(w, r)
		} else {
			mh.getMovie(w, r)
		}
	case http.MethodPost:
		if isCollectionPath {
			mh.postMovies(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	case http.MethodPut:
		if isCollectionPath {
			mh.putMovies(w, r)
		} else {
			mh.putMovie(w, r)
		}
	case http.MethodPatch:
		if isCollectionPath {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		} else {
			mh.patchMovies(w, r)
		}
	case http.MethodDelete:
		if isCollectionPath {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		} else {
			mh.deleteMovie(w, r)
		}
	case http.MethodOptions:
		mh.Options(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func NewMoviesHandler() *MoviesHandler {
	return &MoviesHandler{}
}
