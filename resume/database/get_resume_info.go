package resume_db

import (
	"os"
	"resume-review-api/mongodb"
	"resume-review-api/resume"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetResumeInfo(userId primitive.ObjectID) (int, float64, error) {

	var resumeInfo []resume.ResumesResponse
	filter := bson.D{
		{Key: "user_id", Value: userId},
	}
	err := mongodb.FindMany(os.Getenv("db_name"), "resumes", filter, &resumeInfo)
	if err != nil {
		return 0, 0, err
	}

	// Calculate Count
	count := len(resumeInfo)
	responseTotal := float64(0)

	// Calculate Average Response Time

	if count == 0 {
		return 0, 0, nil
	}

	for _, i := range resumeInfo {
		responseTotal = responseTotal + i.ResponseTime
	}

	averageTime := float64(responseTotal) / float64(count)

	// Send
	return count, averageTime, nil
}

func GetResumeMatched(userId primitive.ObjectID) (int, int, error) {

	var resumeData []resume.DBResumeReview
	filter := bson.D{
		{Key: "user_id", Value: userId},
	}
	err := mongodb.FindMany(os.Getenv("db_name"), "resumes", filter, &resumeData)
	if err != nil {
		return 0, 0, err
	}

	// Calculate no_of resumes
	resumeCount := len(resumeData)
	jobMatch := int(0)

	// Calculate No_Of resumes matched

	if resumeCount == 0 {
		return 0, 0, nil
	}

	for _, resume := range resumeData {
		if resume.ResumeInfo.Score != "" {
			jobMatch++
		}
	}

	// Send
	return resumeCount, jobMatch, nil
}
