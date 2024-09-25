package handlers

import (
	"net/http"
	"notification_service/internal/app/services"
	"notification_service/internal/utils"

	"github.com/labstack/echo/v4"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	NotificationService services.NotificationService
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
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c echo.Context) error {
	notificationDtos, err := h.NotificationService.GetLatestNotifications()
	if err != nil {
		return utils.HandleError(c, err, http.StatusInternalServerError)
	}

	if len(notificationDtos) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, notificationDtos)
}
