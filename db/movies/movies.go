// This package provides methods to interact with the movie database.
package movies

import (
	"database/sql"
	"errors"
	"sync"

	"moviepin/models"
)

type MoviesRepository interface {
	GetMovies() ([]*models.Movie, error)
	GetMovie(id string) (*models.Movie, error)
	AddMovie(movie models.Movie) error
	UpdateMovie(id string, movie models.Movie) error
	DeleteMovie(id string) error
	ReplaceMovies(movies []*models.Movie) error
	GetMovieRating(id string) (*models.MovieReview, error)
}

var (
	// Error returned when movie does not exist.
	ErrNotExists = errors.New("movie does not exist")
)

type Movies struct {
	db *sql.DB
}

func NewMovie(db *sql.DB) *Movies {
	return &Movies{db: db}
}

// Returns slice of all movies present.
func (m Movies) GetMovies() ([]*models.Movie, error) {
	rows, err := m.db.Query("SELECT movie_id, title, release_date, genre, director, description FROM movies;")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	movies := make([]*models.Movie, 0)

	for rows.Next() {
		movie := &models.Movie{}

		if err := rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Genre, &movie.Director, &movie.Description); err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

// Replaces movies collection with passed in collection
func (m Movies) ReplaceMovies(movies []*models.Movie) error {
	tx, err := m.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM movies")

	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, movie := range movies {
		wg.Add(1)

		go func(movie *models.Movie) {
			defer wg.Done()

			_, err := tx.Exec("INSERT INTO movies(movie_id, title, release_date, genre, director, description) VALUES($1, $2, $3, $4, $5, $6);", movie.ID, movie.Title, movie.ReleaseDate, movie.Genre, movie.Director, movie.Description)

			if err != nil {
				tx.Rollback()
				return
			}
		}(movie)
	}

	wg.Wait()

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Returns particular movie.
func (m Movies) GetMovie(id string) (*models.Movie, error) {
	row := m.db.QueryRow("SELECT movie_id, title, release_date, genre, director, description FROM movies WHERE movie_id = $1;", id)

	movie := &models.Movie{}

	if err := row.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Genre, &movie.Director, &movie.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotExists
		}

		return nil, err
	}

	return movie, nil
}

// Adds movie to the database.
func (m Movies) AddMovie(newMovie models.Movie) error {
	if _, err := m.db.Exec("INSERT INTO movies(movie_id, title, release_date, genre, director, description) VALUES($1, $2, $3, $4, $5, $6);", newMovie.ID, newMovie.Title, newMovie.ReleaseDate, newMovie.Genre, newMovie.Director, newMovie.Description); err != nil {
		return err
	}

	return nil
}

// Deletes a movie from the database.
func (m Movies) DeleteMovie(id string) error {
	result, err := m.db.Exec("DELETE FROM movies WHERE movie_id=$1;", id)

	if err != nil {
		return err
	}

	num, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if num == 0 {
		return ErrNotExists
	}

	return nil
}

// Updates a movie in the database.
func (m Movies) UpdateMovie(id string, movie models.Movie) error {
	result, err := m.db.Exec("UPDATE movies SET movie_id=$1, title=$2, release_date=$3, genre=$4, director=$5, description=$6 WHERE movie_id=$7;", movie.ID, movie.Title, movie.ReleaseDate, movie.Genre, movie.Director, movie.Description, id)

	if err != nil {
		return err
	}

	num, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if num == 0 {
		return ErrNotExists
	}

	return nil
}

// Returns movie details along with its rating.
func (m Movies) GetMovieRating(id string) (*models.MovieReview, error) {
	// Take a average of all the ratings for a movie.
	row := m.db.QueryRow(`SELECT m.movie_id, m.title, m.release_date, m.genre, m.director, m.description, TRUNC(ROUND(AVG(r.rating)) / 2, 1) FROM movies m LEFT JOIN reviews r ON m.movie_id=r.movie_id WHERE m.movie_id=$1 GROUP BY m.movie_id;`, id)

	mr := &models.MovieReview{}

	if err := row.Scan(&mr.ID, &mr.Title, &mr.ReleaseDate, &mr.Genre, &mr.Director, &mr.Description, &mr.Rating); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotExists
		}

		return nil, err
	}

	return mr, nil
}
