package resume_db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"resume-review-api/mongodb"
	"resume-review-api/resume"
)

func (s *ResumeDBService) GetResumeInfo(userId primitive.ObjectID) (int, float64, error) {

	var resumeInfo []resume.ResumesResponse
	filter := bson.D{
		{"user_id", userId},
	}
	err := mongodb.FindMany(s.serverSettings.DBName, "resumes", filter, &resumeInfo)
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
