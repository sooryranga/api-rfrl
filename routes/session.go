package routes

import (
	"crypto/rsa"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterSessionRoutes session routes
func RegisterSessionRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PublicKey, sessionUseCase tutorme.SessionUseCase) {

	sessionViews := views.SessionView{SessionUseCase: sessionUseCase}

	sessionR := e.Group("/session")
	sessionR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	sessionR.POST("/", sessionViews.CreateSessionEndpoint)
	sessionR.PUT("/:id", sessionViews.UpdateSessionEndpoint)
	sessionR.DELETE("/:id", sessionViews.DeleteSessionEndpoint)
	sessionR.GET("/:id", sessionViews.GetSessionEndpoint)

	sessionEventR := e.Group("/session/:session-id/event")
	sessionEventR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	sessionEventR.POST("/", sessionViews.CreateSessionEventEndpoint)
	sessionEventR.GET("/:id", sessionViews.GetSessionEventEndpoint)

}
