package routes

import (
	"crypto/rsa"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterClientRoutes register client routes
func RegisterClientRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PublicKey, clientUseCase tutorme.ClientUseCase) {
	clientView := views.ClientView{
		ClientUseCase: clientUseCase,
	}
	r := e.Group("/client")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	validate.RegisterStructValidation(views.ClientPayloadValidation, views.ClientPayload{})
	r.POST("/", clientView.CreateClientEndpoint)
	r.PUT("/:id/", clientView.UpdateClientEndpoint)
	r.GET("/:id/", clientView.GetClientEndpoint)

	r.POST("/:clientID/verify-email/", clientView.VerifyEmail)
	r.PUT("/:clientID/verify-email/", clientView.VerifyEmailPassCode)
	r.GET("/:clientID/verify-email/", clientView.GetVerificationEmails)
	r.DELETE("/:clientID/verify-email/", clientView.DeleteVerifyEmail)

	r.PUT("/:clientID/education/", clientView.CreateEducation)

	tutoringClientR := e.Group("/clients")
	tutoringClientR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	tutoringClientR.GET("/", clientView.GetClientsEndpoint)

	clientEventsR := e.Group("/client/:clientID/events")
	clientEventsR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))
	clientEventsR.GET("/", clientView.GetClientEventsEndpoint)
}
