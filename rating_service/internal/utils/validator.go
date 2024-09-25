package utils

type BaseValidator interface {
	Validate(i interface{}) error
}
