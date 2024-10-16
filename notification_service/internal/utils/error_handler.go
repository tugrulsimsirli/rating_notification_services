package utils

import (
	"github.com/labstack/echo/v4"
)

type ErrorModel struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Errors  []error `json:"errors"`
}

// HandleError is a centralized error handler that returns a consistent error response
func HandleError(c echo.Context, err error, statusCode int) error {
	r := ErrorModel{
		Status:  statusCode,
		Message: err.Error(),
		Errors:  nil,
	}

	return c.JSON(statusCode, r)
}
