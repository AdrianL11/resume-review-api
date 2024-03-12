package authentication_db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"resume-review-api/mongodb"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
)

func ResetPassword(token string, password string) error {

	// Check If Token is Valid
	if !ForgotPasswordValidateToken(token) {
		return errors.New("invalid token")
	}

	// Token Valid, Get User ID from Token
	userId, err := mongodb.GetUserIdByForgotPasswordToken(token)
	if err != nil {
		return err
	}

	// Check if Active User
	profile, err := mongodb.GetProfilebyUserId(userId)
	if err != nil {
		return err
	}

	if profile.ActiveUser == false {
		return errors.New("inactive user")
	}

	// Update Profile
	hashedPassword, _ := util.HashPassword(password)
	filter := bson.D{{"_id", userId}}
	update := bson.D{{"password", hashedPassword}}
	err = mongodb.UpdateOne(resume_ai_env.GetSettingsForEnv().DBName, "users", filter, update)
	if err != nil {
		return err
	}

	// Update Forgot Password Active
	filter = bson.D{{"token", token}}
	update = bson.D{{"is_active", false}}
	err = mongodb.UpdateOne(resume_ai_env.GetSettingsForEnv().DBName, "forgot_passwords", filter, update)
	if err != nil {
		return err
	}

	// Done
	return nil
}
