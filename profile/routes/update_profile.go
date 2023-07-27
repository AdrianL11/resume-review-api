package profile_routes

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"resume-review-api/mongodb"
	profile_util "resume-review-api/profile/util"
)

type UpdateProfileDetails struct {
	FirstName    string `json:"first_name" validate:"omitempty"`
	LastName     string `json:"last_name" validate:"omitempty"`
	Country      string `json:"country" validate:"omitempty"`
	ProfileImage string `json:"profile_image" validate:"omitempty"`
}

func UpdateProfile(c echo.Context) error {

	// Create Update Profile Details
	var updateProfileDetails UpdateProfileDetails
	if err := c.Bind(&updateProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(updateProfileDetails); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validated, Update
	profile, err := mongodb.GetProfileBySession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	filter := bson.D{{"_id", profile.ID}}
	update := bson.D{}

	if updateProfileDetails.FirstName != "" {
		update = append(update, bson.E{"first_name", updateProfileDetails.FirstName})
	}

	if updateProfileDetails.LastName != "" {
		update = append(update, bson.E{"last_name", updateProfileDetails.LastName})
	}

	if updateProfileDetails.Country != "" {
		update = append(update, bson.E{"country", updateProfileDetails.Country})
	}

	image, err := profile_util.GetImageCDNURL(updateProfileDetails.ProfileImage)
	if err != nil {
		update = append(update, bson.E{"profile_image", ""})
	} else {
		update = append(update, bson.E{"profile_image", image})
	}

	err = mongodb.UpdateOne(os.Getenv("db_name"), "users", filter, update)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Done
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
