package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"moviepin/db/movies"
	"moviepin/mocks"
	"moviepin/models"
)

func TestGetMovies(t *testing.T) {
	t.Run("get movies", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMoviesError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovies(rr, req)

		assertStatusCode(t, rr.Code, http.StatusOK)

		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatal(err)
		}

		var movies []*models.Movie

		if err := json.Unmarshal(body, &movies); err != nil {
			t.Fatal(err)
		}

		if len(movies) != 1 {
			t.Errorf("wrong number of movies, got %v want %v", len(movies), 1)
		}

		assertMovie(t, movies[0], &mocks.Movie)
	})

	t.Run("get movies error", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMoviesError = errors.New("error")

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovies(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})
}

func TestGetMovie(t *testing.T) {
	t.Run("get movie", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusOK)

		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatal(err)
		}

		var movie *models.Movie

		if err := json.Unmarshal(body, &movie); err != nil {
			t.Fatal(err)
		}

		assertMovie(t, movie, &mocks.Movie)
	})

	t.Run("get movie wrong path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = errors.New("error")

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("get movie not found", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = movies.ErrNotExists

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNotFound)
	})

	t.Run("get movie error", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = errors.New("error")

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})
}

func TestGetMovieRating(t *testing.T) {
	t.Run("get movie rating", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = nil
		repo.GetMovieRatingError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000?rating=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovieRating(rr, req)

		assertStatusCode(t, rr.Code, http.StatusOK)

		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatal(err)
		}

		var movieReview models.MovieReview

		if err := json.Unmarshal(body, &movieReview); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(&movieReview, &mocks.MovieReview) {
			t.Errorf("wrong movie review, got %v want %v", &movieReview, &mocks.MovieReview)
		}
	})

	t.Run("get movie rating wrong path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/1?rating=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovieRating(rr, req)

		assertStatusCode(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("get movie not found", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = movies.ErrNotExists

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000?rating=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovieRating(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNotFound)
	})

	t.Run("get movie error", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = errors.New("error")

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000?rating=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovieRating(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})

	t.Run("get movie rating error", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieRatingError = errors.New("error")

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000?rating=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.getMovieRating(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})
}

func TestPostMovies(t *testing.T) {
	t.Run("post movie", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.AddMovieError = nil

		handler := NewMoviesHandler(repo)

		body := moviesRequestBody(t, []*models.Movie{&mocks.Movie})

		req, err := http.NewRequest("POST", "/movies", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.postMovies(rr, req)

		assertStatusCode(t, rr.Code, http.StatusCreated)
	})

	t.Run("post movie error all fail", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.AddMovieError = errors.New("error")

		handler := NewMoviesHandler(repo)

		body := moviesRequestBody(t, []*models.Movie{&mocks.Movie})

		req, err := http.NewRequest("POST", "/movies", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.postMovies(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})
}

func TestDeleteMovie(t *testing.T) {
	t.Run("delete movie", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.DeleteMovieError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("DELETE", "/movies/550e8400-e29b-41d4-a716-446655440000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.deleteMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("delete movie wrong path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.DeleteMovieError = errors.New("error")

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("DELETE", "/movies/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.deleteMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("delete movie not found", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.DeleteMovieError = movies.ErrNotExists

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("DELETE", "/movies/550e8400-e29b-41d4-a716-446655440000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.deleteMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNotFound)
	})

	t.Run("delete movie error", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.DeleteMovieError = errors.New("error")

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("DELETE", "/movies/550e8400-e29b-41d4-a716-446655440000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.deleteMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})
}

func TestPutMovie(t *testing.T) {
	t.Run("put movie", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = nil

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PUT", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.putMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("put movie wrong path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PUT", "/movies/1", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.putMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("put movie not found", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = movies.ErrNotExists

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PUT", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.putMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNotFound)
	})

	t.Run("put movie error", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = errors.New("error")

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PUT", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)

		}

		rr := httptest.NewRecorder()

		handler.putMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})
}

func TestPutMovies(t *testing.T) {
	t.Run("put movies", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.ReplaceMoviesError = nil

		handler := NewMoviesHandler(repo)

		body := moviesRequestBody(t, []*models.Movie{&mocks.Movie})

		req, err := http.NewRequest("PUT", "/movies", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.putMovies(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("put movies error", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.ReplaceMoviesError = errors.New("error")

		handler := NewMoviesHandler(repo)

		body := moviesRequestBody(t, []*models.Movie{&mocks.Movie})

		req, err := http.NewRequest("PUT", "/movies", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.putMovies(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})
}

func TestPatchMovie(t *testing.T) {
	t.Run("patch movie updating whole movie", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = nil

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PATCH", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.patchMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("patch movie updating partial movie", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = nil

		handler := NewMoviesHandler(repo)

		movie := make(map[string]interface{})
		movie["title"] = "updated title"

		movieJSON, err := json.Marshal(movie)
		if err != nil {
			t.Fatal(err)
		}

		body := bytes.NewBuffer(movieJSON)

		req, err := http.NewRequest("PATCH", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.patchMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("patch movie wrong path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PATCH", "/movies/1", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.patchMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("patch movie not found", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = movies.ErrNotExists

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PATCH", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.patchMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNotFound)
	})

	t.Run("patch movie error", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = errors.New("error")

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PATCH", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.patchMovie(rr, req)

		assertStatusCode(t, rr.Code, http.StatusInternalServerError)
	})
}

func TestOptions(t *testing.T) {
	t.Run("options", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("OPTIONS", "/movies", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.Options(rr, req)

		if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, PATCH, DELETE, OPTIONS" {
			t.Errorf("wrong Access-Control-Allow-Methods, got %v want %v", rr.Header().Get("Access-Control-Allow-Methods"), "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if got, want := rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type"; got != want {
			t.Errorf("wrong Access-Control-Allow-Headers, got %v want %v", rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
		}

		if got, want := rr.Header().Get("Access-Control-Allow-Origin"), "*"; got != want {
			t.Errorf("wrong Access-Control-Allow-Origin, got %v want %v", rr.Header().Get("Access-Control-Allow-Origin"), "*")
		}

		if got, want := rr.Header().Get("Access-Control-Max-Age"), "86400"; got != want {
			t.Errorf("wrong Access-Control-Max-Age, got %v want %v", rr.Header().Get("Access-Control-Max-Age"), "86400")
		}

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})
}

func TestServeHTTP(t *testing.T) {
	t.Run("get movie path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusOK)
	})

	t.Run("get movies path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMoviesError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusOK)
	})

	t.Run("get movie rating path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieRatingError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("GET", "/movies/550e8400-e29b-41d4-a716-446655440000?rating=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusOK)
	})

	t.Run("post movies path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.AddMovieError = nil

		handler := NewMoviesHandler(repo)

		body := moviesRequestBody(t, []*models.Movie{&mocks.Movie})

		req, err := http.NewRequest("POST", "/movies", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusCreated)
	})

	t.Run("post movie path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()

		handler := NewMoviesHandler(repo)

		body := moviesRequestBody(t, []*models.Movie{&mocks.Movie})

		req, err := http.NewRequest("POST", "/movie", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusMethodNotAllowed)
	})

	t.Run("patch movie path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = nil

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PATCH", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("patch movies path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = nil

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PATCH", "/movies", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusMethodNotAllowed)
	})

	t.Run("put movie path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.UpdateMovieError = nil

		handler := NewMoviesHandler(repo)

		body := movieRequestBody(t, &mocks.Movie)

		req, err := http.NewRequest("PUT", "/movies/550e8400-e29b-41d4-a716-446655440000", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("put movies path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.ReplaceMoviesError = nil

		handler := NewMoviesHandler(repo)

		body := moviesRequestBody(t, []*models.Movie{&mocks.Movie})

		req, err := http.NewRequest("PUT", "/movies", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("delete movie", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.GetMovieError = nil
		repo.DeleteMovieError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("DELETE", "/movies/550e8400-e29b-41d4-a716-446655440000", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("delete movies path", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()
		repo.DeleteMovieError = nil

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest(http.MethodDelete, "/movies", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusMethodNotAllowed)
	})

	t.Run("options", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("OPTIONS", "/movies", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusNoContent)
	})

	t.Run("unknown method", func(t *testing.T) {
		repo := mocks.NewMoviesRepository()

		handler := NewMoviesHandler(repo)

		req, err := http.NewRequest("UNKNOWN", "/movies", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assertStatusCode(t, rr.Code, http.StatusMethodNotAllowed)
	})
}

func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("wrong status code, got %d want %d", got, want)
	}
}

func assertMovie(t *testing.T, got, want *models.Movie) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("wrong movie, got %v want %v", got, want)
	}
}

func moviesRequestBody(t *testing.T, movie []*models.Movie) io.Reader {
	t.Helper()

	movieJSON, err := json.Marshal(movie)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(movieJSON)
}

func movieRequestBody(t *testing.T, movie *models.Movie) io.Reader {
	t.Helper()

	movieJSON, err := json.Marshal(movie)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(movieJSON)
}
