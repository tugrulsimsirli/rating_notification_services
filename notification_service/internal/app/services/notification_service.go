package services

import (
	"encoding/json"
	"log"
	"notification_service/config"
	"notification_service/internal/models/dto"
	"time"
)

type NotificationService struct {
	RabbitMQService RabbitMQServiceInterface
	Config          *config.Config
}

func (s *NotificationService) GetLatestNotifications() ([]dto.NotificationDto, error) {
	var notificationDtos []dto.NotificationDto

	// Consume messages from RabbitMQ
	channel, msgs, err := s.RabbitMQService.CreateChannel(s.Config.RabbitMQ.QueueName)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := s.RabbitMQService.CloseChannel(channel)
		if err != nil {
			log.Printf("Failed to close RabbitMQ channel: %s", err)
		}
	}()

	timeout := time.After(2 * time.Millisecond)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return notificationDtos, nil
			}
			log.Printf("Received a message: %s", msg.Body)

			var notificationDto dto.NotificationDto
			err = json.Unmarshal(msg.Body, &notificationDto)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}

			notificationDtos = append(notificationDtos, notificationDto)

		case <-timeout:
			goto done
		}
	}

done:
	log.Println(notificationDtos)
	return notificationDtos, nil
}
