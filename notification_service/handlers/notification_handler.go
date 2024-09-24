package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"notification_service/models/dto"
	"notification_service/services"
	_ "notification_service/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tugrulsimsirli/rabbitmq"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	NotificationService services.NotificationService
	RabbitMQService     *rabbitmq.RabbitMQService
}

// GetNotifications retrieves the latest notifications
// @Summary Get latest notifications
// @Description Get the latest notifications about ratings
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 200 {array} dto.NotificationDto
// @Failure 204 "No Content, no new notifications"
// @Failure 400 {object} utils.ErrorModel
// @Router /notifications/{providerID} [get]
func (h *NotificationHandler) GetNotifications(c echo.Context) error {
	var notificationDtos []dto.NotificationDto

	// Yeni bir RabbitMQ kanalı oluştur
	channel, msgs, err := h.RabbitMQService.CreateChannel("notification_queue")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	defer func() {
		err := h.RabbitMQService.CloseChannel(channel)
		if err != nil {
			log.Printf("Failed to close RabbitMQ channel: %s", err)
		}
	}()

	// Mesajları dinlemek için bir timeout ekleyelim (örneğin 2 saniye)
	timeout := time.After(200 * time.Millisecond)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				// Mesaj kanalı kapalıysa çık
				return c.NoContent(http.StatusNoContent)
			}
			log.Printf("Received a message: %s", msg.Body)

			// Mesajı DTO modeline çevir
			var notificationDto dto.NotificationDto
			err = json.Unmarshal(msg.Body, &notificationDto)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}

			notificationDtos = append(notificationDtos, notificationDto)

		case <-timeout:
			// Zaman aşımı dolduysa döngüden çık
			goto done
		}
	}

done:
	// Eğer ilgili provider için mesaj yoksa 204 No Content dön
	if len(notificationDtos) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	// Mesajları DTO formatında dön
	return c.JSON(http.StatusOK, notificationDtos)
}
