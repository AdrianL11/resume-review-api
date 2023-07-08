package session_db

import (
	"errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"resume-review-api/mongodb"
	"time"
)

func ValidateSession(c echo.Context) error {

	// Check if Session is Valid from Cookie
	sess, err := session.Get("_resumereview-tpl", c)
	if err != nil {
		return err
	}

	if sess.Values["email_address"] == nil || sess.Values["session_id"] == nil {
		return errors.New("invalid session")
	}

	emailAddress := sess.Values["email_address"].(string)
	sessionId := sess.Values["session_id"].(string)

	if emailAddress == "" || sessionId == "" {
		return errors.New("invalid session")
	}

	// Cooke is Good, Check DB if Session Exists
	var sessionInformation mongodb.Session
	filter := bson.D{{"session_id", sessionId}}
	err = mongodb.FindOne("resume_reviewer", "sessions", filter, &sessionInformation)
	if err != nil {
		return err
	}

	// Session found in DB, Is it Active
	if sessionInformation.Active == false {
		return errors.New("invalid session")
	}

	// Is it Expired?
	expiration := time.Unix(int64(sessionInformation.Expiration.T), 0).UTC()
	now := time.Now().UTC()

	if now.After(expiration) {
		return errors.New("invalid session")
	}

	// Is User Active?
	profile, err := mongodb.GetProfilebyUserId(sessionInformation.UserId)
	if err != nil {
		return err
	}

	if profile.ActiveUser == false {
		return errors.New("invalid session")
	}

	// Update Last Seen
	err = UpdateLastSeen(c)
	if err != nil {
		return err
	}

	// Not Expired, Good
	return nil
}
