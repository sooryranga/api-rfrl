package routes

import (
	"crypto/rsa"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterReportClient register report client routes
func RegisterReportClient(e *echo.Echo, validate *validator.Validate, publicKey *rsa.PublicKey, reportClientUseCase tutorme.ReportClientUseCase) {
	views := views.ReportClientView{ReportClientUseCase: reportClientUseCase}

	reportR := e.Group("/report")
	reportR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))
	reportR.POST("/", views.CreateReport)
	reportR.DELETE("/", views.DeleteReport)
	reportR.GET("/", views.GetReports)
}
