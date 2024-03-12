package mongodb

import (
	"errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/util/resume_ai_env"
)

func GetUserIdByEmail(emailAddress string) (primitive.ObjectID, error) {

	var profile Profile
	filter := bson.D{{"email_address", emailAddress}, {"active_user", true}}

	err := FindOne(resume_ai_env.GetSettingsForEnv().DBName, "users", filter, &profile)
	if err != nil {
		return profile.ID, err
	}

	return profile.ID, nil
}

func GetUserIdByForgotPasswordToken(token string) (primitive.ObjectID, error) {

	var forgotPasswordDetails ForgotPassword
	filter := bson.D{{"token", token}}

	err := FindOne(resume_ai_env.GetSettingsForEnv().DBName, "forgot_passwords", filter, &forgotPasswordDetails)
	if err != nil {
		return forgotPasswordDetails.UserId, err
	}

	return forgotPasswordDetails.UserId, err
}

func GetProfilebyEmailAddress(emailAddress string) (Profile, error) {

	var profile Profile
	filter := bson.D{{"email_address", emailAddress}, {"active_user", true}}

	err := FindOne(resume_ai_env.GetSettingsForEnv().DBName, "users", filter, &profile)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func GetProfilebyUserId(id primitive.ObjectID) (Profile, error) {

	var profile Profile
	filter := bson.D{{"_id", id}}

	err := FindOne(resume_ai_env.GetSettingsForEnv().DBName, "users", filter, &profile)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func GetProfileBySession(c echo.Context) (Profile, error) {

	var profile Profile

	// Get Session Data
	sess, err := session.Get(resume_ai_env.GetSettingsForEnv().SessionCookieName, c)
	if err != nil {
		return profile, nil
	}

	// Grab Email from Session
	emailAddress := sess.Values["email_address"].(string)
	if emailAddress == "" {
		return profile, errors.New("no valid session")
	}

	// Get Profile
	profile, err = GetProfilebyEmailAddress(emailAddress)
	if err != nil {
		return profile, nil
	}

	return profile, nil
}
