package utils

import (
	"fmt"
	"rating_service/internal/utils/validators"
	"reflect"
)

type ValidatorFactory struct{}

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
