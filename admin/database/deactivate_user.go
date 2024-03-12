package admin_db

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"resume-review-api/resume_ai_middleware"
	"time"
)

func (s *ResumeAIAdminDBService) DeactivateUser(c echo.Context, objId primitive.ObjectID) error {
	viewerProfile, ok := c.Get(resume_ai_middleware.UserSessionProfile).(mongodb.Profile)
	if !ok {
		return echo.ErrInternalServerError
	}

	// Find User
	filter := bson.D{{"_id", objId}}
	update := bson.D{
		{"active_user", false},
		{"deactivation_date", time.Now().UTC()},
		{"deactivated_by", viewerProfile.ID},
	}
	err := mongodb.UpdateOne(s.serverSettings.DBName, "users", filter, update)
	if err != nil {
		return err
	}

	return nil
}
