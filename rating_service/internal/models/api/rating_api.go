package api

import "github.com/google/uuid"

type RatingRequest struct {
	ProviderID uuid.UUID `json:"provider_id" example:"d1dcc1d6-f1ee-49e8-acf6-d7a5de7f4bea" validate:"required,uuid"` // UUID and required validation
	Rating     float64   `json:"rating" validate:"required,gte=1,lte=5"`                                              // Rating should be between 1 and 5
}
