package profile_routes

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
	"resume-review-api/util/resume_ai_env"
)

func GetProfile(c echo.Context) error {

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Session is Valid, Return Profile
	sess, err := session.Get(resume_ai_env.GetSettingsForEnv().SessionCookieName, c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	profile, err := mongodb.GetProfilebyEmailAddress(sess.Values["email_address"].(string))
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, profile)
}
