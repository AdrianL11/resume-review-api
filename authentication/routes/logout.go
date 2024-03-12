package authentication_routes

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *AuthRouteHandler) Logout(c echo.Context) error {

	// Validated, Revoke Session
	sess, err := session.Get(h.serverSettings.SessionCookieName, c)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "success",
		})
	}

	emailAddress := sess.Values["session_id"].(string)
	if emailAddress == "" {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "success",
		})
	}

	err = h.authService.InvalidateSession(emailAddress)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "success",
		})
	}

	sess.Options.MaxAge = -1
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// Return Success
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
