package admin_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"resume-review-api/mongodb"
)

func GetUsers() ([]mongodb.Profile, error) {

	var profiles []mongodb.Profile
	var err error

	// Grab All Profiles
	filter := bson.D{}
	err = mongodb.FindMany(os.Getenv("db_name"), "users", filter, &profiles)

	return profiles, err
}
