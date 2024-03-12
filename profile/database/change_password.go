package profile_db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
)

func ChangePassword(id primitive.ObjectID, oldPassword string, newPassword string) error {

	// Get Profile from ID
	var oldPasswordDB mongodb.ChangePassword
	filter := bson.D{{"_id", id}}
	err := mongodb.FindOne(resume_ai_env.GetSettingsForEnv().DBName, "users", filter, &oldPasswordDB)
	if err != nil {
		return err
	}

	// Does Old Password Match?
	if !util.CheckPasswordHash(oldPassword, oldPasswordDB.Password) {
		return errors.New("incorrect old password")
	}

	// Matches, Lets Change
	hashedPassword, _ := util.HashPassword(newPassword)
	update := mongodb.ChangePassword{
		Password: hashedPassword,
	}
	err = mongodb.UpdateOne(resume_ai_env.GetSettingsForEnv().DBName, "users", filter, update)
	if err != nil {
		return err
	}

	return nil
}
