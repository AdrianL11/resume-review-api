package authentication_routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *AuthRouteHandler) LoggedIn(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"role": h.authService.LoggedIn(c),
	})
}
