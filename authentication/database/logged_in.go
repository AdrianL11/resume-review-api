package authentication_db

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"resume-review-api/mongodb"
	"time"
)

type SessionValidation struct {
	Expiration primitive.Timestamp     `bson:"expiration" json:"expiration"`
	UserInfo   []SessionValidationUser `bson:"result" json:"result"`
}

type SessionValidationUser struct {
	ActiveUser bool   `bson:"active_user" json:"active_user"`
	Role       string `bson:"role" json:"role"`
}

func LoggedIn(c echo.Context) string {

	// Check if Session is Valid from Cookie
	sess, err := session.Get("_resumereview-tpl", c)
	if err != nil {
		return ""
	}

	if sess.Values["email_address"] == nil || sess.Values["session_id"] == nil {
		return ""
	}

	emailAddress := sess.Values["email_address"].(string)
	sessionId := sess.Values["session_id"].(string)

	if emailAddress == "" || sessionId == "" {
		return ""
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
	err = mongodb.Aggregate("resume_reviewer", "sessions", mongo.Pipeline{matchStage, lookupStage}, &mongoUser)
	if err != nil {
		return ""
	}

	// Any Results?
	if len(mongoUser) <= 0 {
		return ""
	}

	// Is Session Expired?
	expiration := time.Unix(int64(mongoUser[0].Expiration.T), 0).UTC()
	now := time.Now().UTC()

	if now.After(expiration) {
		return ""
	}

	// Is User Active?
	if !mongoUser[0].UserInfo[0].ActiveUser {
		return ""
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
		return ""
	}

	// Done
	return mongoUser[0].UserInfo[0].Role
}
