// Mock for the movies repository interface.
package mocks

import (
	"moviepin/models"
	"time"

	"github.com/google/uuid"
)

var (
	Movie = models.Movie{
		ID:          uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
		Title:       "The Shawshank Redemption",
		ReleaseDate: time.Date(1994, time.September, 23, 0, 0, 0, 0, time.UTC),
		Genre:       "Drama",
		Director:    "Frank Darabont",
		Description: "Prisoners",
	}

	MovieReview = models.MovieReview{
		ID:          Movie.ID,
		Title:       Movie.Title,
		ReleaseDate: Movie.ReleaseDate,
		Genre:       Movie.Genre,
		Director:    Movie.Director,
		Description: Movie.Description,
		Rating:      3.5,
	}
)

// MoviesRepository is a mock for the movies repository interface.
type MoviesRepository struct {
	GetMoviesError      error
	GetMovieError       error
	AddMovieError       error
	UpdateMovieError    error
	DeleteMovieError    error
	ReplaceMoviesError  error
	GetMovieRatingError error
}

// NewMoviesRepository returns a new instance of the movies repository mock.
func NewMoviesRepository() MoviesRepository {
	return MoviesRepository{}
}

// GetMovie returns a movie by its id.
func (m MoviesRepository) GetMovie(id string) (*models.Movie, error) {
	if m.GetMovieError != nil {
		return nil, m.GetMovieError
	}

	return &Movie, nil
}

// GetMovies returns a slice of all movies present.
func (m MoviesRepository) GetMovies() ([]*models.Movie, error) {
	if m.GetMoviesError != nil {
		return nil, m.GetMoviesError
	}

	return []*models.Movie{&Movie}, nil
}

// AddMovie adds a movie to the database.
func (m MoviesRepository) AddMovie(movie models.Movie) error {
	if m.AddMovieError != nil {
		return m.AddMovieError
	}

	return nil
}

// UpdateMovie updates a movie in the database.
func (m MoviesRepository) UpdateMovie(id string, movie models.Movie) error {
	if m.UpdateMovieError != nil {
		return m.UpdateMovieError
	}

	return nil
}

// DeleteMovie deletes a movie from the database.
func (m MoviesRepository) DeleteMovie(id string) error {
	if m.DeleteMovieError != nil {
		return m.DeleteMovieError
	}

	return nil
}

// ReplaceMovies replaces all movies in the database.
func (m MoviesRepository) ReplaceMovies(movies []*models.Movie) error {
	if m.ReplaceMoviesError != nil {
		return m.ReplaceMoviesError
	}

	return nil
}

// GetMovieRating returns a movie rating by its id.
func (m MoviesRepository) GetMovieRating(id string) (*models.MovieReview, error) {
	if m.GetMovieRatingError != nil {
		return nil, m.GetMovieRatingError
	}

	return &MovieReview, nil
}
