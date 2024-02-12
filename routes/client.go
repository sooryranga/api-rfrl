package routes

import (
	"crypto/rsa"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/Arun4rangan/api-rfrl/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterClientRoutes register client routes
func RegisterClientRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PublicKey, clientUseCase rfrl.ClientUseCase) {
	clientView := views.ClientView{
		ClientUseCase: clientUseCase,
	}
	r := e.Group("/client")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
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

	r.PUT("/:clientID/wanting-company-referral/", clientView.CreateWantingReferralCompany)
	r.GET("/:clientID/wanting-company-referral/", clientView.GetWantingReferralCompany)

	clientsR := e.Group("/clients")
	clientsR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	clientsR.GET("/", clientView.GetClientsEndpoint)

	clientEventsR := e.Group("/client/:clientID/events")
	clientEventsR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))
	clientEventsR.GET("/", clientView.GetClientEventsEndpoint)
}
