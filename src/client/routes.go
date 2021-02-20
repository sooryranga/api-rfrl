package client

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	algorithmRS256 = "RS256"
)

// Register auth routes
func (h *Handler) Register(e *echo.Echo, validate *validator.Validate) {
	r := e.Group("/client")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    h.key,
		SigningMethod: algorithmRS256,
	}))

	validate.RegisterStructValidation(clientPayloadValidation, ClientPayload{})
	r.POST("/", h.CreateClientEndpoint)
	r.PUT("/:id", h.UpdateClientEndpoint)
	r.GET("/:id", h.GetClientEndpoint)
}
