package authentication_db

import (
	aws_ses "resume-review-api/aws-ses"
	profile_db "resume-review-api/profile/database"
	"resume-review-api/util/resume_ai_env"
)

type ResumeAIAuthDBService struct {
	serverSettings   resume_ai_env.ServerSettings
	profileDBService *profile_db.ProfileDBService
	emailService     *aws_ses.EmailService
}

func NewResumeAIAuthDBService(
	serverSettings resume_ai_env.ServerSettings,
	profileDBService *profile_db.ProfileDBService,
	emailService *aws_ses.EmailService,
) *ResumeAIAuthDBService {
	return &ResumeAIAuthDBService{
		serverSettings:   serverSettings,
		profileDBService: profileDBService,
		emailService:     emailService,
	}
}
