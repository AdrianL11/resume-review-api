package admin_db

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"time"
)

func DeactivateUser(c echo.Context, objId primitive.ObjectID) error {

	// Current User
	profile, err := mongodb.GetProfileBySession(c)
	if err != nil {
		return err
	}

	// Find User
	filter := bson.D{{"_id", objId}}
	update := bson.D{
		{"active_user", false},
		{"deactivation_date", primitive.Timestamp{T: uint32(time.Now().UTC().Unix())}},
		{"deactivated_by", profile.ID},
	}
	err = mongodb.UpdateOne("resume_reviewer", "users", filter, update)
	if err != nil {
		return err
	}

	return nil
}
