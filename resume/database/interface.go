package resume_db

import (
	"resume-review-api/util/resume_ai_env"
)

type ResumeDBService struct {
	serverSettings resume_ai_env.ServerSettings
}

func NewResumeDBService(serverSettings resume_ai_env.ServerSettings) *ResumeDBService {
	return &ResumeDBService{
		serverSettings: serverSettings,
	}
}
