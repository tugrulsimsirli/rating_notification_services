package utils

import (
	"fmt"
	"rating_service/internal/utils/validators"
	"reflect"
)

// ValidatorFactory provides a method to get the appropriate validator for each model
type ValidatorFactory struct{}

// NewValidatorFactory returns a new instance of ValidatorFactory
func NewValidatorFactory() *ValidatorFactory {
	return &ValidatorFactory{}
}

// GetValidator returns the appropriate validator for the given model type
func (vf *ValidatorFactory) GetValidator(model interface{}) (BaseValidator, error) {
	modelType := reflect.TypeOf(model)

	switch modelType {
	case reflect.TypeOf(validators.RatingRequestValidator{}):
		return validators.NewRatingRequestValidator(), nil
	default:
		return nil, fmt.Errorf("no validator found for model type: %s", modelType)
	}
}
