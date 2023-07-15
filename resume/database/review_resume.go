package resume_db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"resume-review-api/mongodb"
	"resume-review-api/resume"
	"time"
)

func InsertResumeReview(userId primitive.ObjectID, resumeObj resume.JSONObject, responseTime float64) error {

	var doc = resume.DBResumeReview{
		UserId:       userId,
		ResponseTime: responseTime,
		CreationDate: time.Now().UTC(),
		ExpiresAt:    time.Now().UTC().Add(time.Hour * 24 * 30), // 30 Day Expiration
		ResumeInfo:   resumeObj,
	}

	_, err := mongodb.NewDocument(os.Getenv("db_name"), "resumes", doc)
	if err != nil {
		return err
	}

	// Inserted
	return nil
}
