package session_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"resume-review-api/mongodb"
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
	err := mongodb.UpdateMany(os.Getenv("db_name"), "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}
