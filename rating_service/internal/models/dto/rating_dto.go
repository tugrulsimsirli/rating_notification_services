package dto

import "github.com/google/uuid"

type RatingDto struct {
	ProviderID uuid.UUID `json:"provider_id"`
	Rating     float64   `json:"rating"`
}

type AverageRatingResponseDto struct {
	AverageRating float64 `json:"average_rating"`
}
