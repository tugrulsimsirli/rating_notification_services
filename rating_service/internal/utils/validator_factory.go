package utils

import (
	"fmt"
	"rating_service/internal/utils/validators"
)

// ValidatorFactory provides a method to get the appropriate validator for each model
type ValidatorFactory struct{}

// NewValidatorFactory returns a new instance of ValidatorFactory
func NewValidatorFactory() *ValidatorFactory {
	return &ValidatorFactory{}
}

// GetValidator returns the appropriate validator for the given model
func (vf *ValidatorFactory) GetValidator(modelName string) (BaseValidator, error) {
	switch modelName {
	case "RatingRequest":
		return validators.NewRatingRequestValidator(), nil // Ensure it's properly initialized
	default:
		return nil, fmt.Errorf("no validator found for model: %s", modelName)
	}
}
