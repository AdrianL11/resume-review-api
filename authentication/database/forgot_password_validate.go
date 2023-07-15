package authentication_db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"resume-review-api/mongodb"
	"time"
)

func ForgotPasswordValidateToken(token string) bool {

	// Check If Token Exists
	var forgotPasswordStruct mongodb.ForgotPassword
	filter := bson.D{{"token", token}}

	if err := mongodb.FindOne(os.Getenv("db_name"), "forgot_passwords", filter, &forgotPasswordStruct); err != nil {
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
