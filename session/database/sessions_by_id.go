package session_db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"resume-review-api/mongodb"
)

func GetSessionsById(id primitive.ObjectID) ([]mongodb.Session, error) {

	// Lookup Sessions
	var sessions []mongodb.Session
	var err error

	filter := bson.D{
		{"user_id", id},
		{"is_active", true},
	}
	err = mongodb.FindMany(os.Getenv("db_name"), "sessions", filter, &sessions)
	if err != nil || len(sessions) <= 0 {
		return sessions, errors.New("no results")
	}

	return sessions, err
}
