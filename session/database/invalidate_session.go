package session_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"resume-review-api/mongodb"
	"time"
)

func InvalidateSession(sessionId string) error {

	// Invalidate
	filter := bson.D{{"session_id", sessionId}}
	update := bson.D{
		{"is_active", false},
		{"revocation_date", time.Now().UTC()},
	}
	err := mongodb.UpdateOne(os.Getenv("db_name"), "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}
