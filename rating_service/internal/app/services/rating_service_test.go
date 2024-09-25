package services_test

import (
	"fmt"
	"rating_service/internal/app/services"
	"rating_service/internal/db/repositories"
	"rating_service/internal/models/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRabbitMQService struct {
	mock.Mock
}

func (m *MockRabbitMQService) Publish(message string) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestAddRating(t *testing.T) {
	db, sql_mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock oluşturulamadı: %s", err)
	}
	defer db.Close()

	mockRabbitMQService := new(MockRabbitMQService)

	ratingRepo := repositories.RatingRepository{DB: db}
	ratingService := services.RatingService{
		RatingRepository: ratingRepo,
		RabbitMQService:  mockRabbitMQService,
	}

	ratingDTO := dto.RatingDto{
		ProviderID: uuid.New(),
		Rating:     4.5,
	}

	sql_mock.ExpectBegin()
	sql_mock.ExpectExec("INSERT INTO ratings").WithArgs(sqlmock.AnyArg(), float64(4.5)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql_mock.ExpectCommit()

	mockRabbitMQService.On("Publish", mock.AnythingOfType("string")).Return(nil)

	err = ratingService.AddRating(ratingDTO)
	assert.NoError(t, err)

	mockRabbitMQService.AssertExpectations(t)

	err = sql_mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCalculateAverageRating(t *testing.T) {
	db, sql_mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock oluşturulamadı: %s", err)
	}
	defer db.Close()

	ratingRepo := repositories.RatingRepository{DB: db}
	ratingService := services.RatingService{
		RatingRepository: ratingRepo,
	}

	providerID := uuid.New()
	expectedAvgRating := dto.AverageRatingResponseDto{
		AverageRating: 4.5,
	}

	rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(float64(4.5))
	sql_mock.ExpectQuery("SELECT COALESCE\\(AVG\\(rating\\), 0\\) FROM ratings WHERE provider_id = \\$1").
		WithArgs(providerID).
		WillReturnRows(rows)

	avgRating, err := ratingService.CalculateAverageRating(providerID)
	assert.NoError(t, err)

	assert.Equal(t, expectedAvgRating, avgRating)

	err = sql_mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddRating_DBError(t *testing.T) {
	db, sql_mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock oluşturulamadı: %s", err)
	}
	defer db.Close()

	mockRabbitMQService := new(MockRabbitMQService)

	ratingRepo := repositories.RatingRepository{DB: db}
	ratingService := services.RatingService{
		RatingRepository: ratingRepo,
		RabbitMQService:  mockRabbitMQService,
	}

	ratingDTO := dto.RatingDto{
		ProviderID: uuid.New(),
		Rating:     4.5,
	}

	sql_mock.ExpectBegin()
	sql_mock.ExpectExec("INSERT INTO ratings").WithArgs(sqlmock.AnyArg(), float64(4.5)).
		WillReturnError(fmt.Errorf("DB error"))
	sql_mock.ExpectRollback()

	err = ratingService.AddRating(ratingDTO)
	assert.Error(t, err)

	err = sql_mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddRating_RabbitMQError(t *testing.T) {
	db, sql_mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock oluşturulamadı: %s", err)
	}
	defer db.Close()

	mockRabbitMQService := new(MockRabbitMQService)

	ratingRepo := repositories.RatingRepository{DB: db}
	ratingService := services.RatingService{
		RatingRepository: ratingRepo,
		RabbitMQService:  mockRabbitMQService,
	}

	ratingDTO := dto.RatingDto{
		ProviderID: uuid.New(),
		Rating:     4.5,
	}

	sql_mock.ExpectBegin()
	sql_mock.ExpectExec("INSERT INTO ratings").WithArgs(sqlmock.AnyArg(), 4.5).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql_mock.ExpectCommit()

	mockRabbitMQService.On("Publish", mock.AnythingOfType("string")).Return(fmt.Errorf("RabbitMQ error"))

	err = ratingService.AddRating(ratingDTO)
	assert.Error(t, err)

	mockRabbitMQService.AssertExpectations(t)
	err = sql_mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
