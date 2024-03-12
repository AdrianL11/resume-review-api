package profile_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
)

func (s *ProfileDBService) GetUserIdByEmail(emailAddress string) (primitive.ObjectID, error) {

	var profile mongodb.Profile
	filter := bson.D{{"email_address", emailAddress}, {"active_user", true}}

	err := mongodb.FindOne(s.serverSettings.DBName, "users", filter, &profile)
	if err != nil {
		return profile.ID, err
	}

	return profile.ID, nil
}

func (s *ProfileDBService) GetUserIdByForgotPasswordToken(token string) (primitive.ObjectID, error) {

	var forgotPasswordDetails mongodb.ForgotPassword
	filter := bson.D{{"token", token}}

	err := mongodb.FindOne(s.serverSettings.DBName, "forgot_passwords", filter, &forgotPasswordDetails)
	if err != nil {
		return forgotPasswordDetails.UserId, err
	}

	return forgotPasswordDetails.UserId, err
}

func (s *ProfileDBService) GetProfileByEmailAddress(emailAddress string) (mongodb.Profile, error) {

	var profile mongodb.Profile
	filter := bson.D{{"email_address", emailAddress}, {"active_user", true}}

	err := mongodb.FindOne(s.serverSettings.DBName, "users", filter, &profile)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (s *ProfileDBService) GetProfileByUserId(id primitive.ObjectID) (mongodb.Profile, error) {

	var profile mongodb.Profile
	filter := bson.D{{"_id", id}}

	err := mongodb.FindOne(s.serverSettings.DBName, "users", filter, &profile)
	if err != nil {
		return profile, err
	}

	return profile, nil
}
