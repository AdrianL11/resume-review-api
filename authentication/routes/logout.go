package authentication_routes

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	session_db "resume-review-api/session/database"
)

func Logout(c echo.Context) error {

	// Validate Session
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Validated, Revoke Session
	sess, err := session.Get("_resumereview-tpl", c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	emailAddress := sess.Values["session_id"].(string)
	if emailAddress == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	err = session_db.InvalidateSession(emailAddress)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
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
