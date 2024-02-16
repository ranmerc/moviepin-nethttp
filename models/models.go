package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required"`
	ReleaseDate time.Time `json:"release_date" validate:"required"`
	Genre       string    `json:"genre" validate:"required"`
	Director    string    `json:"director" validate:"required"`
	Description string    `json:"description" validate:"required"`
}

type MovieReview struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required"`
	ReleaseDate time.Time `json:"release_date" validate:"required"`
	Genre       string    `json:"genre" validate:"required"`
	Director    string    `json:"director" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Rating      float32   `json:"rating" validate:"required,lte=5,gte=0"`
}

type Review struct {
	ID         uuid.UUID `json:"id" validate:"required,uuid"`
	UserID     uuid.UUID `json:"user_id" validate:"required,uuid"`
	MovieID    uuid.UUID `json:"movie_id" validate:"required,uuid"`
	Rating     float32   `json:"rating" validate:"required,lte=5,gte=0"`
	ReviewText string    `json:"review_text" validate:"required,lte=500"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
	UpdatedAt  time.Time `json:"updated_at" validate:"required"`
}
