package routes

import (
	"crypto/rsa"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/Arun4rangan/api-rfrl/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterCompanyRoutes register auth routes
func RegisterCompanyRoutes(e *echo.Echo, validate *validator.Validate, publicKey *rsa.PublicKey, companyUseCase rfrl.CompanyUseCase) {
	views := views.CompanyView{CompanyUseCase: companyUseCase}

	companyR := e.Group("/company")
	companyR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))
	companyR.POST("/", views.CreateCompanyView)
	companyR.GET("/", views.GetCompanies)
	companyR.PUT("/:id/", views.UpdateCompanyView)
	companyR.GET("/:id/", views.GetCompany)

	companyEmailR := e.Group("/company-email")
	companyEmailR.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    publicKey,
		SigningMethod: rfrl.AlgorithmRS256,
		Claims:        &rfrl.JWTClaims{},
	}))
	companyEmailR.PUT("/", views.UpdateCompanyEmailView)
	companyEmailR.GET("/", views.GetCompanyEmailsView)
	companyEmailR.GET("/:companyEmail/", views.GetCompanyEmailView)
}
