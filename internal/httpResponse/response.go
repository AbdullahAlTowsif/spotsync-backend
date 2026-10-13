package httpResponse

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type SuccessBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, SuccessBody{Success: true, Message: message, Data: data})
}

func Created(c *echo.Context, message string, data any) error {
	return c.JSON(http.StatusCreated, SuccessBody{Success: true, Message: message, Data: data})
}
