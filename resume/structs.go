package resume

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type JSONObject struct {
	Name             string       `json:"name"`
	Email            string       `json:"email"`
	Phone            string       `json:"phone"`
	Address          string       `json:"address"`
	Summary          string       `json:"summary"`
	Skills           []string     `json:"skills"`
	Experiences      []Experience `json:"experiences"`
	Educations       []Education  `json:"educations"`
	Score            string       `json:"score"`
	ScoreReason      string       `json:"scoreReason"`
	RecruiterSummary string       `json:"recruiter_summary"`
}

type Experience struct {
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Dates    string   `json:"dates"`
	JobTitle string   `json:"job_title"`
	Duties   []string `json:"duties"`
}

type Education struct {
	SchoolName     string `json:"school_name"`
	Location       string `json:"location"`
	Type           string `json:"type"`
	GraduationYear string `json:"graduation"`
}

type DBResumeReview struct {
	UserId       primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreationDate time.Time          `bson:"creation_date" json:"creation_date"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expires_at"`
	ResumeInfo   JSONObject         `bson:"resume_info" json:"resume_info"`
	ResponseTime float64            `bson:"response_time" json:"response_time"`
}

type ResumesResponse struct {
	ResponseTime float64 `bson:"response_time" json:"response_time"`
}

type ReturnResponseResumeInfo struct {
	Count        int     `bson:"count" json:"count"`
	ResponseTime float64 `bson:"response_time" json:"response_time"`
	PerDay       float64 `bson:"per_day" json:"per_day"`
}
