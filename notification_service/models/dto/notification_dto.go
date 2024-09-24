
package dto

import "github.com/google/uuid"

type NotificationDto struct {
    Id         uuid.UUID `json:"id"`
    ProviderID uuid.UUID `json:"provider_id"`
    Message    string    `json:"message"`
}
