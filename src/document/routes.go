package document

import (
	"github.com/Arun4rangan/api-tutorme/src/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	algorithmRS256 = "RS256"
)

// Register auth routes
func (h *Handler) Register(e *echo.Echo) {
	r := e.Group("/document")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    h.key,
		SigningMethod: algorithmRS256,
		Claims:        auth.JWTClaims{},
	}))

	r.POST("/", h.CreateDocumentEndpoint)
	r.PUT("/:id", h.UpdateDocumentEndpoint)
	r.DELETE("/:id", h.DeleteDocumentEndpoint)
	r.GET("/:id", h.GetDocumentEndpoint)

	r2 := e.Group("/document-order")
	r2.POST("/", h.CreateDocumentOrderEndpoint)
	r2.PUT("/", h.UpdateDocumentOrderEndpoint)
	r2.GET("/", h.GetDocumentOrderEndpoint)
}
