package profile_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"resume-review-api/mongodb"
	profile_util "resume-review-api/profile/util"
	"resume-review-api/util"
)

func SetProfile(token string, password string, firstName string, lastName string, country string, profileImage string) error {

	// Create ObjectID
	obj, _ := primitive.ObjectIDFromHex(token)

	// Insert
	hashedPassword, _ := util.HashPassword(password)
	filter := bson.D{{"_id", obj}}

	image, err := profile_util.GetImageCDNURL(profileImage)
	if err != nil {
		profileImage = ""
	} else {
		profileImage = image
	}

	doc := mongodb.InsertProfile{
		Password:     hashedPassword,
		FirstLogin:   false,
		FirstName:    firstName,
		LastName:     lastName,
		Country:      country,
		ProfileImage: profileImage,
	}

	err = mongodb.UpdateOne(os.Getenv("db_name"), "users", filter, doc)
	if err != nil {
		return err
	}

	// Done
	return nil
}
