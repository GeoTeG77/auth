package router

import (
	"auth/internal/api/handlers"
	"auth/internal/service"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, tokenManager *service.TokenManager) *echo.Echo {
	handler := &handlers.Handler{
		TokenManager: tokenManager,
	}

	e.GET("/api/v1/token/guid", handler.GetToken)
	e.GET("/api/v1/token/refresh", handler.UpdateToken)

	return e
}
