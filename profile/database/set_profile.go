package profile_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"resume-review-api/profile/cdn"
	"resume-review-api/util"
)

func (s *ProfileDBService) SetProfile(token string, password string, firstName string, lastName string, country string, profileImage string) error {

	// Create ObjectID
	obj, _ := primitive.ObjectIDFromHex(token)

	// Insert
	hashedPassword, _ := util.HashPassword(password)
	filter := bson.D{{"_id", obj}}

	image, err := cdn.GetImageCDNURL(profileImage)
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

	err = mongodb.UpdateOne(s.serverSettings.DBName, "users", filter, doc)
	if err != nil {
		return err
	}

	// Done
	return nil
}
