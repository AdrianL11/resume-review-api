package authentication_db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"resume-review-api/mongodb"
	"resume-review-api/util"
)

func CheckLogin(emailAddress string, password string) error {

	// Check If Allowed
	var profile mongodb.ChangePassword
	filter := bson.D{{"email_address", emailAddress}, {"active_user", true}}
	err := mongodb.FindOne(os.Getenv("db_name"), "users", filter, &profile)
	if err != nil {
		return err
	}

	// Check Password
	allowed := util.CheckPasswordHash(password, profile.Password)
	if !allowed {
		return errors.New("unauthorized access")
	}

	return nil
}
