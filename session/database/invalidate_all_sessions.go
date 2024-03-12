package session_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"resume-review-api/util/resume_ai_env"
	"time"
)

func InvalidateAllSessions(userId string) error {

	// Create Primitive ObjectID
	userIdObj, _ := primitive.ObjectIDFromHex(userId)

	// Invalidate
	filter := bson.D{{"user_id", userIdObj}}
	update := bson.D{
		{"is_active", false},
		{"revocation_date", time.Now().UTC()},
	}
	err := mongodb.UpdateMany(resume_ai_env.GetSettingsForEnv().DBName, "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}
