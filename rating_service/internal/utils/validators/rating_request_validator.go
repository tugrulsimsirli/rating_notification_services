package validators

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// RatingRequestValidator validates RatingRequest structs
type RatingRequestValidator struct {
	validate *validator.Validate
}

// NewRatingRequestValidator creates a new instance of RatingRequestValidator
func NewRatingRequestValidator() *RatingRequestValidator {
	return &RatingRequestValidator{
		validate: validator.New(), // Properly initialize the validator
	}
}

// Validate performs validation on a RatingRequest struct
func (v *RatingRequestValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

// HandleValidationError provides custom error handling for RatingRequest validation errors
func (v *RatingRequestValidator) HandleValidationError(c echo.Context, err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// Format validation errors specific to RatingRequest
		errors := make(map[string]string)
		for _, validationError := range validationErrors {
			switch validationError.Tag() {
			case "required":
				errors[validationError.Field()] = validationError.Field() + " is required"
			case "gte":
				errors[validationError.Field()] = validationError.Field() + " must be greater than or equal to " + validationError.Param()
			case "lte":
				errors[validationError.Field()] = validationError.Field() + " must be less than or equal to " + validationError.Param()
			default:
				errors[validationError.Field()] = "Invalid value for " + validationError.Field()
			}
		}

		// Return a structured error response
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Fallback for non-validation errors
	return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
}
