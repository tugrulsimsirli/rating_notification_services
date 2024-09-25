
package data

import "github.com/google/uuid"

type Rating struct {
    ID         uuid.UUID `json:"id"`
    ProviderID uuid.UUID `json:"provider_id"`
    Rating     float64   `json:"rating"`
    CreatedAt  string    `json:"created_at"`
    IsDeleted  bool      `json:"is_deleted"`
}
