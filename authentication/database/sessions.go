package authentication_db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"time"
)

func (s *ResumeAIAuthDBService) InvalidateAllSessions(userId string) error {

	// Create Primitive ObjectID
	userIdObj, _ := primitive.ObjectIDFromHex(userId)

	// Invalidate
	filter := bson.D{{"user_id", userIdObj}}
	update := bson.D{
		{"is_active", false},
		{"revocation_date", time.Now().UTC()},
	}
	err := mongodb.UpdateMany(s.serverSettings.DBName, "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *ResumeAIAuthDBService) InvalidateSession(sessionId string) error {

	// Invalidate
	filter := bson.D{{"session_id", sessionId}}
	update := bson.D{
		{"is_active", false},
		{"revocation_date", time.Now().UTC()},
	}
	err := mongodb.UpdateOne(s.serverSettings.DBName, "sessions", filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *ResumeAIAuthDBService) GetSessionsById(id primitive.ObjectID) ([]mongodb.Session, error) {
	// Lookup Sessions
	var sessions []mongodb.Session
	var err error

	filter := bson.D{
		{"user_id", id},
		{"is_active", true},
	}
	err = mongodb.FindMany(s.serverSettings.DBName, "sessions", filter, &sessions)
	if err != nil || len(sessions) <= 0 {
		return sessions, errors.New("no results")
	}

	return sessions, err
}
