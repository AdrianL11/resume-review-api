package admin_db

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"resume-review-api/mongodb"
	profile_util "resume-review-api/profile/util"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
)

func UpdateProfile(id string, email string, firstName string, lastName string, country string, profileImage string, password string, role string) error {

	obj, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", obj}}
	update := bson.D{}

	if email != "" {
		update = append(update, bson.E{"email_address", email})
	}

	if firstName != "" {
		update = append(update, bson.E{"first_name", firstName})
	}

	if lastName != "" {
		update = append(update, bson.E{"last_name", lastName})
	}

	if profileImage != "" {

		image, err := profile_util.GetImageCDNURL(profileImage)
		if err != nil {
			update = append(update, bson.E{"profile_image", ""})
		} else {
			update = append(update, bson.E{"profile_image", image})
		}
	}

	if role != "" {
		update = append(update, bson.E{"role", role})
	}

	if password != "" {
		hashedPass, _ := util.HashPassword(password)
		update = append(update, bson.E{"password", hashedPass})
	}

	if country != "" {
		update = append(update, bson.E{"country", country})
	}

	err := mongodb.UpdateOne(resume_ai_env.GetSettingsForEnv().DBName, "users", filter, update)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
