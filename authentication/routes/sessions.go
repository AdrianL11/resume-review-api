package authentication_routes

import (
	"errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"resume-review-api/mongodb"
	"time"
)

func (h *AuthRouteHandler) CreateSession(c echo.Context) error {

	// Get Session
	sess, err := session.Get(h.serverSettings.SessionCookieName, c)
	if err != nil {
		return err
	}

	// Get Session Email Address
	sessionEmail := sess.Values["email_address"].(string)
	if sessionEmail == "" {
		return errors.New("invalid session key")
	}

	// Get User Profile from Email Address
	var profile mongodb.Profile
	filter := bson.D{{"email_address", sessionEmail}}
	err = mongodb.FindOne(h.serverSettings.DBName, "users", filter, &profile)
	if err != nil {
		return err
	}

	// Is User Active
	if profile.ActiveUser == false {
		return errors.New("unauthorized access")
	}

	// Add Session to Collection
	sessionId := sess.Values["session_id"].(string)
	if sessionId == "" {
		return errors.New("invalid session key")
	}

	doc := mongodb.Session{
		ID:         sessionId,
		LoggedInIP: c.RealIP(),
		LastSeenIP: c.RealIP(),
		Expiration: time.Now().UTC().Add(time.Hour * 24 * 14),
		LastSeen:   time.Now().UTC(),
		Active:     true,
		UserAgent:  c.Request().UserAgent(),
		UserId:     profile.ID,
	}

	if _, err := mongodb.NewDocument(h.serverSettings.DBName, "sessions", doc); err != nil {
		return err
	}

	return nil
}
