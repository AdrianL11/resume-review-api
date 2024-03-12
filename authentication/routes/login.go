package authentication_routes

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"resume-review-api/util/resume_ai_env"
)

type LoginDetails struct {
	Email    string `json:"email_address" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthRouteHandler) Login(c echo.Context) error {

	// Create Login Details
	var loginDetails LoginDetails
	if err := c.Bind(&loginDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(loginDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check DB for Login Details
	err := h.authService.CheckLogin(loginDetails.Email, loginDetails.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Found, Create Session
	sess, err := session.Get(h.serverSettings.SessionCookieName, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if resume_ai_env.IsProd() {
		sess.Options = &sessions.Options{
			MaxAge:   3600 * 24 * 14, // 14 Days
			Domain:   h.serverSettings.SessionCookieDomain,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		}
	} else {
		sess.Options = &sessions.Options{
			MaxAge:   3600 * 24 * 14, // 14 Days
			Domain:   h.serverSettings.SessionCookieDomain,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		}
	}

	sess.Values["email_address"] = loginDetails.Email
	sess.Values["session_id"] = uuid.New().String()
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Add Session to Database
	err = h.CreateSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Return Success
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
