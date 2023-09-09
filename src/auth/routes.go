package auth

import "github.com/labstack/echo/v4"

// Register auth routes
func (handler *Handler) Register(e *echo.Echo) {
	e.POST("/login", handler.Login)
}
