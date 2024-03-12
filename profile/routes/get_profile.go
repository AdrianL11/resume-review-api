package profile_routes

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *ProfileRouteHandler) GetProfile(c echo.Context) error {
	// Session is Valid, Return Profile
	sess, err := session.Get(h.serverSettings.SessionCookieName, c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	profile, err := h.profileService.GetProfileByEmailAddress(sess.Values["email_address"].(string))
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, profile)
}
