
package services

import (
    "notification_service/models/dto"
    "notification_service/repositories"
    "github.com/google/uuid"
)

type NotificationService struct {
    NotificationRepository repositories.NotificationRepository
}

func (s *NotificationService) GetLatestNotifications(providerID uuid.UUID) ([]dto.NotificationDto, error) {
    return s.NotificationRepository.GetNotifications(providerID)
}
