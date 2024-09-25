package services

import (
	"encoding/json"
	"fmt"
	"rating_service/internal/db/repositories"
	"rating_service/internal/models/dto"

	"github.com/google/uuid"
)

type RatingService struct {
	RatingRepository repositories.RatingRepository
	RabbitMQService  RabbitMQServiceInterface
}

func (s *RatingService) AddRating(ratingDto dto.RatingDto) error {
	err := s.RatingRepository.InsertRating(ratingDto)
	if err != nil {
		return err
	}

	// Send message to RabbitMQ
	message, err := json.Marshal(dto.NotificationDto{
		Id:         uuid.New(),
		ProviderID: ratingDto.ProviderID,
		Message:    fmt.Sprintf("Provider %s got a rating of %.1f", ratingDto.ProviderID, ratingDto.Rating),
	})
	if err != nil {
		return err
	}

	err = s.RabbitMQService.Publish(string(message))
	if err != nil {
		return err
	}

	return nil
}

func (s *RatingService) CalculateAverageRating(providerID uuid.UUID) (dto.AverageRatingResponseDto, error) {
	return s.RatingRepository.GetAverageRating(providerID)
}
