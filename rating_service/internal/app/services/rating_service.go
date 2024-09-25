package services

import (
	"rating_service/internal/db/repositories"
	"rating_service/internal/models/dto"

	"github.com/google/uuid"
)

type RatingService struct {
	RatingRepository repositories.RatingRepository
}

func (s *RatingService) AddRating(ratingDto dto.RatingDto) error {
	return s.RatingRepository.InsertRating(ratingDto)
}

func (s *RatingService) CalculateAverageRating(providerID uuid.UUID) (dto.AverageRatingResponseDto, error) {
	return s.RatingRepository.GetAverageRating(providerID)
}
