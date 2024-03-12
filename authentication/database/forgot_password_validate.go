package authentication_db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"resume-review-api/mongodb"
	"resume-review-api/util/resume_ai_env"
	"time"
)

func ForgotPasswordValidateToken(token string) bool {

	// Check If Token Exists
	var forgotPasswordStruct mongodb.ForgotPassword
	filter := bson.D{{"token", token}}

	if err := mongodb.FindOne(resume_ai_env.GetSettingsForEnv().DBName, "forgot_passwords", filter, &forgotPasswordStruct); err != nil {
		fmt.Println(err.Error())
		return false
	}

	// Check if Active User
	profile, err := mongodb.GetProfilebyUserId(forgotPasswordStruct.UserId)
	if err != nil {
		return false
	}

	if profile.ActiveUser == false {
		return false
	}

	// Is it Active?
	if !forgotPasswordStruct.Active {
		fmt.Println(forgotPasswordStruct.Active)
		return false
	}

	// Is it Expired
	expiration := forgotPasswordStruct.Expiration
	now := time.Now().UTC()

	if now.After(expiration) {
		return false
	}

	return true
}
