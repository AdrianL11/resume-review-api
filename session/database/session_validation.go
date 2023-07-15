package session_db

import (
	"errors"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"resume-review-api/mongodb"
	"time"
)

type SessionValidation struct {
	Expiration time.Time               `bson:"expiration" json:"expiration"`
	UserInfo   []SessionValidationUser `bson:"result" json:"result"`
}

type SessionValidationUser struct {
	ActiveUser bool `bson:"active_user" json:"active_user"`
}

func ValidateSession(c echo.Context) error {

	// Check if Session is Valid from Cookie
	sess, err := session.Get(os.Getenv("session_name"), c)
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

	// Cooke is Good, Grab Information from Session
	var mongoUser []SessionValidation

	matchStage := bson.D{
		{"$match",
			bson.D{
				{"session_id", sessionId},
				{"is_active", true},
			},
		},
	}

	lookupStage := bson.D{
		{"$lookup",
			bson.D{
				{"from", "users"},
				{"localField", "user_id"},
				{"foreignField", "_id"},
				{"as", "result"},
			},
		},
	}

	// Aggregate Groups Created, Lets Look Up
	err = mongodb.Aggregate(os.Getenv("db_name"), "sessions", mongo.Pipeline{matchStage, lookupStage}, &mongoUser)
	if err != nil {
		fmt.Printf("Aggregate Error: %s\n", err.Error())
		return err
	}

	// Any Results?
	if len(mongoUser) <= 0 {
		return errors.New("invalid session")
	}

	// Is Session Expired?
	expiration := mongoUser[0].Expiration
	now := time.Now().UTC()

	if now.After(expiration) {
		return errors.New("invalid session")
	}

	// Is User Active?
	if !mongoUser[0].UserInfo[0].ActiveUser {
		return errors.New("not an active user")
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

	// Not Expired, Good
	return nil
}
