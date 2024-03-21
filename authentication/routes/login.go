package authentication_routes

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	authentication_db "resume-review-api/authentication/database"
	session_db "resume-review-api/session/database"
)

type LoginDetails struct {
	Email    string `json:"email_address" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c echo.Context) error {

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
	err := authentication_db.CheckLogin(loginDetails.Email, loginDetails.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Found, Create Session
	sess, err := session.Get(os.Getenv("session_name"), c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sess.Options = &sessions.Options{
		MaxAge:   3600 * 24 * 14, // 14 Days
		Domain:   "vdart.ai",
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	sess.Values["email_address"] = loginDetails.Email
	sess.Values["session_id"] = uuid.New().String()
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Add Session to Database
	err = session_db.CreateSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Return Success
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
