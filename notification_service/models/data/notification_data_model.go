package data

import "github.com/google/uuid"

type Notification struct {
	ID         uuid.UUID `json:"id"`
	ProviderID uuid.UUID `json:"provider_id"`
	Message    string    `json:"message"`
	CreatedAt  string    `json:"created_at"`
	IsDeleted  bool      `json:"is_deleted"`
}
