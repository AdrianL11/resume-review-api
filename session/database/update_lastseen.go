package session_db

import (
	"errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"time"
)

func UpdateLastSeen(c echo.Context) error {

	// Session ID
	sess, err := session.Get("_resumereview-tpl", c)
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
		{"last_seen", primitive.Timestamp{T: uint32(time.Now().UTC().Unix())}},
		{"user_agent", c.Request().UserAgent()},
	}
	err = mongodb.UpdateOne("resume_reviewer", "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}
