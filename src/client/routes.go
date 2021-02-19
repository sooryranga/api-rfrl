package client

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Register auth routes
func (h *Handler) Register(e *echo.Echo, validate *validator.Validate) {
	r := e.Group("/client")
	r.Use(middleware.JWT(h.key))

	validate.RegisterStructValidation(clientPayloadValidation, ClientPayload{})
	r.POST("/client", h.CreateClientEndpoint)
	r.PUT("/client/:id", h.UpdateClientEndpoint)
	r.GET("/client/:id", h.GetClientEndpoint)
}
