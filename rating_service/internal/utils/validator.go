package utils

// BaseValidator defines the structure for all model-specific validators
type BaseValidator interface {
	Validate(i interface{}) error
}
