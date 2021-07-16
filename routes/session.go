package routes

import (
	"crypto/rsa"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/Arun4rangan/api-rfrl/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterSessionRoutes session routes
func RegisterSessionRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PublicKey, sessionUseCase rfrl.SessionUseCase, tutorReviewUseCase rfrl.TutorReviewUseCase) {

	sessionViews := views.SessionView{SessionUseCase: sessionUseCase, TutorReviewUseCase: tutorReviewUseCase}

	sessionR := e.Group("/session")
	sessionR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	sessionR.POST("/", sessionViews.CreateSessionEndpoint)
	sessionR.PUT("/:id/", sessionViews.UpdateSessionEndpoint)
	sessionR.DELETE("/:id/", sessionViews.DeleteSessionEndpoint)
	sessionR.GET("/:id/", sessionViews.GetSessionEndpoint)

	sessionEventR := e.Group("/session/:sessionID/event")
	sessionEventR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	sessionEventR.POST("/", sessionViews.CreateSessionEventEndpoint)
	sessionEventR.GET("/:id/", sessionViews.GetSessionEventEndpoint)
	sessionEventR.GET("/", sessionViews.GetSessionRelatedEventsEndpoint)

	clientActionOnEventR := e.Group("/session/:sessionId/book")
	clientActionOnEventR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	clientActionOnEventR.POST("/", sessionViews.CreateClientActionOnSessionEvent)

	sessionConferenceR := e.Group("/session/:sessionID/conference")
	sessionConferenceR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	sessionConferenceR.GET("/", sessionViews.GetSessionConferenceIDEndpoint)

	sessionsR := e.Group("/sessions/")
	sessionsR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))

	sessionsR.GET("", sessionViews.GetSessionsEndpoint)
}
