package admin_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"resume-review-api/mongodb"
	"resume-review-api/util/resume_ai_env"
)

func GetUsers() ([]mongodb.Profile, error) {

	var profiles []mongodb.Profile
	var err error

	// Grab All Profiles
	filter := bson.D{}
	err = mongodb.FindMany(resume_ai_env.GetSettingsForEnv().DBName, "users", filter, &profiles)

	return profiles, err
}
