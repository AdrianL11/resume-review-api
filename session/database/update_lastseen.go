package session_db

import (
	"errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"resume-review-api/mongodb"
	"time"
)

func UpdateLastSeen(c echo.Context) error {

	// Session ID
	sess, err := session.Get(os.Getenv("session_name"), c)
	if err != nil {
		return err
	}

	sessionId := sess.Values["session_id"].(string)
	if sessionId == "" {
		return errors.New("invalid session key")
	}

	// Update Session
	filter := bson.D{{"session_id", sessionId}}
	update := bson.D{
		{"lastseen_ip", c.RealIP()},
		{"last_seen", time.Now().UTC()},
		{"user_agent", c.Request().UserAgent()},
	}
	err = mongodb.UpdateOne(os.Getenv("db_name"), "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}
