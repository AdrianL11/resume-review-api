package authentication_routes

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	session_db "resume-review-api/session/database"
	"resume-review-api/util/resume_ai_env"
)

func Logout(c echo.Context) error {

	// Validate Session
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "success",
		})
	}

	// Validated, Revoke Session
	sess, err := session.Get(resume_ai_env.GetSettingsForEnv().SessionCookieName, c)
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

	err = session_db.InvalidateSession(emailAddress)
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
