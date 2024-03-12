package session_db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"resume-review-api/util/resume_ai_env"
)

func GetSessionsById(id primitive.ObjectID) ([]mongodb.Session, error) {

	// Lookup Sessions
	var sessions []mongodb.Session
	var err error

	filter := bson.D{
		{"user_id", id},
		{"is_active", true},
	}
	err = mongodb.FindMany(resume_ai_env.GetSettingsForEnv().DBName, "sessions", filter, &sessions)
	if err != nil || len(sessions) <= 0 {
		return sessions, errors.New("no results")
	}

	return sessions, err
}
