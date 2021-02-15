package auth

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// Register auth routes
func (handler *Handler) Register(e *echo.Echo, validate *validator.Validate) {

	validate.RegisterStructValidation(LoginPayloadValidation, LoginPayload{})
	e.POST("/login", handler.Login)

	validate.RegisterStructValidation(SignUpPayloadValidation, SignUpPayload{})
	e.POST("/signup", handler.Signup)
}
