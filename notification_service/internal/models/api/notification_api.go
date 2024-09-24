package api

import "github.com/google/uuid"

type NotificationResponse struct {
	Id      uuid.UUID `json:"id" example:"d1dcc1d6-f1ee-49e8-acf6-d7a5de7f4bea"`
	Message string    `json:"message"`
}
