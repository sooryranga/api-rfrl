package routes

import (
	"crypto/rsa"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/Arun4rangan/api-rfrl/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterCompanyRoutes register auth routes
func RegisterConferenceRoutes(
	e *echo.Echo,
	publicKey *rsa.PublicKey,
	apiKey string,
	sessionUseCase rfrl.SessionUseCase,
	conferenceUseCase rfrl.ConferenceUseCase,
) {
	views := views.ConferenceView{SessionUseCase: sessionUseCase, ConferenceUseCase: conferenceUseCase}

	conferenceR := e.Group("/conference/:conferenceID")
	conferenceR.GET("/yjs/", views.ConnectToSessionYJSClients)
	conferenceR.GET("/simple-peer/", views.ConnectToSessionSimplePeerClients)

	conferenceSessionR := e.Group("conference-session/:sessionID")
	conferenceSessionR.POST("/code/", views.SubmitCode, middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))
	conferenceSessionR.POST("/code/:ID/", views.SetCodeResult, middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == apiKey, nil
	}))
}
