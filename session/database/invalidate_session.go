package session_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"time"
)

func InvalidateSession(sessionId string) error {

	// Invalidate
	filter := bson.D{{"session_id", sessionId}}
	update := bson.D{
		{"is_active", false},
		{"revocation_date", primitive.Timestamp{T: uint32(time.Now().UTC().Unix())}},
	}
	err := mongodb.UpdateOne("resume_reviewer", "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}
