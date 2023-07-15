package profile_db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"resume-review-api/mongodb"
	"resume-review-api/util"
)

func ChangePassword(id primitive.ObjectID, oldPassword string, newPassword string) error {

	// Get Profile from ID
	var oldPasswordDB mongodb.ChangePassword
	filter := bson.D{{"_id", id}}
	err := mongodb.FindOne(os.Getenv("db_name"), "users", filter, &oldPasswordDB)
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
	err = mongodb.UpdateOne(os.Getenv("db_name"), "users", filter, update)
	if err != nil {
		return err
	}

	return nil
}
