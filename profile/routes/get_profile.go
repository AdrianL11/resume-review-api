package profile_routes

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"resume-review-api/mongodb"
	session_db "resume-review-api/session/database"
)

func GetProfile(c echo.Context) error {

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Session is Valid, Return Profile
	sess, err := session.Get(os.Getenv("session_name"), c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	profile, err := mongodb.GetProfilebyEmailAddress(sess.Values["email_address"].(string))
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, profile)
}
