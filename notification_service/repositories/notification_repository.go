
package repositories

import (
    "database/sql"
    "notification_service/models/dto"
    "github.com/google/uuid"
)

type NotificationRepository struct {
    DB *sql.DB
}

func (r *NotificationRepository) GetNotifications(providerID uuid.UUID) ([]dto.NotificationDto, error) {
    var notifications []dto.NotificationDto
    rows, err := r.DB.Query("SELECT id, message FROM notifications WHERE provider_id = $1 AND is_deleted = false", providerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var notification dto.NotificationDto
        err := rows.Scan(&notification.Id, &notification.Message)
        if err != nil {
            return nil, err
        }
        notifications = append(notifications, notification)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    _, err = r.DB.Exec("UPDATE notifications SET is_deleted = true WHERE provider_id = $1", providerID)
    if err != nil {
        return nil, err
    }

    return notifications, nil
}
