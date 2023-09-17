package user

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Register auth routes
func (h *Handler) Register(e *echo.Echo, validate *validator.Validate) {
	r := e.Group("/user")
	r.Use(middleware.JWT(h.key))

	validate.RegisterStructValidation(userPayloadValidation, UserPayload{})
	r.POST("/user", h.CreateUserEndpoint)
	r.PUT("/user/:id", h.UpdateUserEndpoint)
	r.GET("/user/:id", h.GetUserEndpoint)
}
