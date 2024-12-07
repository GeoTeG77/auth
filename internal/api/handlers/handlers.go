package handlers

import (
	"auth/internal/service"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	TokenManager *service.TokenManager
}

func (h *Handler) GetToken(c echo.Context) error {
	guid := c.Param("guid")
	ip := c.Request().RemoteAddr
	JWT, err := h.TokenManager.GetToken(guid, ip)
	if err != nil {
		slog.Error("Bad Request")
		return c.JSON(http.StatusBadRequest, nil)
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    JWT.AccessToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   os.Getenv("URI"),
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    JWT.RefreshToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   os.Getenv("URI"),
	}

	c.SetCookie(accessCookie)
	c.SetCookie(refreshCookie)

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) UpdateToken(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Error("Cookie Not Found")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Refresh token not found"})
	}
	ip := c.Request().RemoteAddr

	JWT, err := h.TokenManager.CheckToken(refreshToken.String(), ip)

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    JWT.AccessToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   os.Getenv("URI"),
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    JWT.RefreshToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   os.Getenv("URI"),
	}

	c.SetCookie(accessCookie)
	c.SetCookie(refreshCookie)

	return c.JSON(http.StatusOK, nil)
}
