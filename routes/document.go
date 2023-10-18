package routes

import (
	"crypto/rsa"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/Arun4rangan/api-tutorme/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterDocumentRoutes documentt routes
func RegisterDocumentRoutes(e *echo.Echo, validate *validator.Validate, key *rsa.PublicKey, documentUseCase tutorme.DocumentUseCase) {
	documentViews := views.DocumentView{DocumentUseCase: documentUseCase}

	r := e.Group("/document")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))

	r.POST("/", documentViews.CreateDocumentEndpoint)
	r.PUT("/:id/", documentViews.UpdateDocumentEndpoint)
	r.DELETE("/:id/", documentViews.DeleteDocumentEndpoint)
	r.GET("/:id/", documentViews.GetDocumentEndpoint)

	r2 := e.Group("/document-order")
	r2.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: tutorme.AlgorithmRS256,
		Claims:        &tutorme.JWTClaims{},
	}))
	r2.POST("/", documentViews.CreateDocumentOrderEndpoint)
	r2.PUT("/", documentViews.UpdateDocumentOrderEndpoint)
	r2.GET("/", documentViews.GetDocumentOrderEndpoint)
}
