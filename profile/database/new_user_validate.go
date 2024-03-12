package profile_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
)

func (s *ProfileDBService) NewUserValidateToken(token string) bool {

	// Create ObjectID
	obj, _ := primitive.ObjectIDFromHex(token)

	// Check If Token Exists
	var profile mongodb.Profile
	filter := bson.D{{"_id", obj}}

	if err := mongodb.FindOne(s.serverSettings.DBName, "users", filter, &profile); err != nil {
		return false
	}

	// Check if Active User
	if profile.ActiveUser == false {
		return false
	}

	// Is it First Time Logging In?
	if profile.FirstLogin == false {
		return false
	}

	return true
}
