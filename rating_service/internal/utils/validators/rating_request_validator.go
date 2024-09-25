package validators

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RatingRequestValidator struct {
	validate *validator.Validate
}

func NewRatingRequestValidator() *RatingRequestValidator {
	return &RatingRequestValidator{
		validate: validator.New(),
	}
}

func (v *RatingRequestValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func (v *RatingRequestValidator) HandleValidationError(c echo.Context, err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
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

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
}
