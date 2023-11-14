package routes

import (
	"crypto/rsa"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/labstack/echo/v4"
)

// RegisterCompanyRoutes register auth routes
func RegisterConferenceRoutes(e *echo.Echo, publicKey *rsa.PublicKey, sessionUseCase tutorme.SessionUseCase, conferenceUseCase tutorme.ConferenceUseCase) {
	views := views.ConferenceView{SessionUseCase: sessionUseCase, ConferenceUseCase: conferenceUseCase}

	conferenceR := e.Group("/conference/:conferenceID")
	conferenceR.GET("/", views.ConnectToSessionClients)
}
