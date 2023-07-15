package profile_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"resume-review-api/mongodb"
)

func NewUserValidateToken(token string) bool {

	// Create ObjectID
	obj, _ := primitive.ObjectIDFromHex(token)

	// Check If Token Exists
	var profile mongodb.Profile
	filter := bson.D{{"_id", obj}}

	if err := mongodb.FindOne(os.Getenv("db_name"), "users", filter, &profile); err != nil {
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
