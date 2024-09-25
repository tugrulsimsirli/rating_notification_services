package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rating_service/internal/app/services"
	"rating_service/internal/models/api"
	"rating_service/internal/models/dto"
	"rating_service/internal/utils"
	"rating_service/internal/utils/validators"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tugrulsimsirli/rabbitmq"
)

// RatingHandler handles rating-related requests
type RatingHandler struct {
	RatingService   services.RatingService
	RabbitMQService *rabbitmq.RabbitMQService
}

// SubmitRating handles the submission of a rating
// @Summary Submit a rating
// @Description Submit a new rating for a service provider
// @Tags Ratings
// @Accept json
// @Produce json
// @Param rating body api.RatingRequest true "Rating"
// @Success 204
// @Failure 400 {object} utils.ErrorModel
// @Router /submit-rating [post]
func (h *RatingHandler) SubmitRating(c echo.Context) error {
	var ratingReq api.RatingRequest
	if err := c.Bind(&ratingReq); err != nil {
		return utils.HandleError(c, err, http.StatusBadRequest)
	}

	// Initialize the validator factory
	validatorFactory := utils.NewValidatorFactory()

	// Get the specific validator for RatingRequest
	validator, err := validatorFactory.GetValidator("RatingRequest")
	if err != nil {
		return utils.HandleError(c, err, http.StatusInternalServerError)
	}

	// Validate the request
	if err := validator.Validate(&ratingReq); err != nil {
		// Call HandleValidationError from the specific validator
		return validator.(*validators.RatingRequestValidator).HandleValidationError(c, err)
	}

	// Map API request to DTO
	var ratingDTO dto.RatingDto
	utils.Map(&ratingReq, &ratingDTO)

	err = h.RatingService.AddRating(ratingDTO)
	if err != nil {
		return utils.HandleError(c, err, http.StatusInternalServerError)
	}

	// RabbitMQ'ya mesajı JSON formatında gönderelim
	message, err := json.Marshal(dto.NotificationDto{
		Id:         uuid.New(),
		ProviderID: ratingDTO.ProviderID,
		Message:    fmt.Sprintf("Provider %s got a rating of %.1f", ratingDTO.ProviderID, ratingDTO.Rating),
	})
	if err != nil {
		return utils.HandleError(c, err, http.StatusInternalServerError)
	}

	err = h.RabbitMQService.Publish(string(message))
	if err != nil {
		return utils.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

// GetAverageRating handles retrieving the average rating for a service provider
// @Summary Get average rating
// @Description Get the average rating for a specific service provider
// @Tags Ratings
// @Accept  json
// @Produce  json
// @Param   providerID path string true "Provider ID"
// @Success 200 {object} dto.AverageRatingResponseDto
// @Failure 400 {object} utils.ErrorModel
// @Router /average-rating/{providerID} [get]
func (h *RatingHandler) GetAverageRating(c echo.Context) error {
	providerID, err := uuid.Parse(c.Param("providerID"))
	if err != nil {
		return utils.HandleError(c, err, http.StatusBadRequest)
	}

	avgRatingDto, err := h.RatingService.CalculateAverageRating(providerID)
	if err != nil {
		return utils.HandleError(c, err, http.StatusInternalServerError)
	} else if avgRatingDto.AverageRating == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, avgRatingDto)
}
