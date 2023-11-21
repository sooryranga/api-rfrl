package routes

import (
	"crypto/rsa"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterCompanyRoutes register auth routes
func RegisterConferenceRoutes(e *echo.Echo, publicKey *rsa.PublicKey, sessionUseCase tutorme.SessionUseCase, conferenceUseCase tutorme.ConferenceUseCase) {
	views := views.ConferenceView{SessionUseCase: sessionUseCase, ConferenceUseCase: conferenceUseCase}

	conferenceR := e.Group("/conference/:conferenceID")
	conferenceR.GET("/", views.ConnectToSessionClients)

	conferenceSessionR := e.Group("conference-session/:sessionId")
	conferenceSessionR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))
	conferenceSessionR.POST("/code/", views.SubmitCode)
}
