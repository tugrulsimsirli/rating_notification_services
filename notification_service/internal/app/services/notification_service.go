package services

import (
	"notification_service/internal/db/repositories"
	"notification_service/internal/models/dto"

	"github.com/google/uuid"
)

type NotificationService struct {
	NotificationRepository repositories.NotificationRepository
}

func (s *NotificationService) GetLatestNotifications(providerID uuid.UUID) ([]dto.NotificationDto, error) {
	return s.NotificationRepository.GetNotifications(providerID)
}
