package session_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"resume-review-api/mongodb"
	"resume-review-api/util/resume_ai_env"
	"time"
)

func InvalidateSession(sessionId string) error {

	// Invalidate
	filter := bson.D{{"session_id", sessionId}}
	update := bson.D{
		{"is_active", false},
		{"revocation_date", time.Now().UTC()},
	}
	err := mongodb.UpdateOne(resume_ai_env.GetSettingsForEnv().DBName, "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}
