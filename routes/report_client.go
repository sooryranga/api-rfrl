package routes

import (
	"crypto/rsa"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/Arun4rangan/api-rfrl/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterReportClient register report client routes
func RegisterReportClient(e *echo.Echo, validate *validator.Validate, publicKey *rsa.PublicKey, reportClientUseCase rfrl.ReportClientUseCase) {
	views := views.ReportClientView{ReportClientUseCase: reportClientUseCase}

	reportR := e.Group("/report")
	reportR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))
	reportR.POST("/", views.CreateReport)
	reportR.DELETE("/", views.DeleteReport)
	reportR.GET("/", views.GetReports)
}
