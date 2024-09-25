package repositories

import (
	"database/sql"
	"fmt"
	"rating_service/internal/models/dto"

	"github.com/google/uuid"
)

type RatingRepository struct {
	DB *sql.DB
}

func (r *RatingRepository) InsertRating(rating dto.RatingDto) error {
	// Start the transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// First INSERT query inside the transaction
	_, err = tx.Exec("INSERT INTO ratings (provider_id, rating) VALUES ($1, $2);", rating.ProviderID, rating.Rating)
	if err != nil {
		tx.Rollback() // Roll back the transaction if there’s an error
		return err
	}

	// Second INSERT query inside the transaction (example: insert into notifications)
	message := fmt.Sprintf("%s got a rating as %.1f", rating.ProviderID.String(), rating.Rating)
	_, err = tx.Exec("INSERT INTO notifications (provider_id ,message) VALUES ($1, $2);", rating.ProviderID, message)
	if err != nil {
		tx.Rollback() // Roll back the transaction if there’s an error
		return err
	}

	// Commit the transaction if both inserts succeed
	return tx.Commit()
}

func (r *RatingRepository) GetAverageRating(providerID uuid.UUID) (dto.AverageRatingResponseDto, error) {
	var avgRating float64
	err := r.DB.QueryRow("SELECT COALESCE(AVG(rating), 0) FROM ratings WHERE provider_id = $1", providerID).Scan(&avgRating)
	if err != nil {
		return dto.AverageRatingResponseDto{}, err
	}

	// Map the result to AverageRatingResponse DTO
	averageRatingResponse := dto.AverageRatingResponseDto{
		AverageRating: float64(int(avgRating*10)) / 10,
	}

	return averageRatingResponse, nil
}
