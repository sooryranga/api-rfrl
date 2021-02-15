package user

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Register auth routes
func (h *Handler) Register(e *echo.Echo) {
	r := e.Group("/user")
	r.Use(middleware.JWT(h.key))
	e.POST("/login", handler.Login)
	e.POST("/signup", handler.Signup)
}
