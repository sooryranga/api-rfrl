package routes

import (
	"crypto/rsa"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// RegisterAuthRoutes register auth routes
func RegisterAuthRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PrivateKey, authUseCase tutorme.AuthUseCase) {
	authview := views.AuthView{
		AuthUseCases: authUseCase,
		Key:          *key,
	}
	validate.RegisterStructValidation(views.LoginPayloadValidation, views.LoginPayload{})
	e.POST("/login/", authview.Login)

	validate.RegisterStructValidation(views.SignUpPayloadValidation, views.SignUpPayload{})
	e.POST("/signup/", authview.Signup)
}
