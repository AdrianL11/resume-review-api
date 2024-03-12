package admin_db

import (
	aws_ses "resume-review-api/aws-ses"
	profile_db "resume-review-api/profile/database"
	"resume-review-api/util/resume_ai_env"
)

type ResumeAIAdminDBService struct {
	serverSettings   resume_ai_env.ServerSettings
	profileDBService *profile_db.ProfileDBService
	emailService     *aws_ses.EmailService
}

func NewResumeAIAdminDBService(
	serverSettings resume_ai_env.ServerSettings,
	profileDBService *profile_db.ProfileDBService,
	emailService *aws_ses.EmailService,
) *ResumeAIAdminDBService {
	return &ResumeAIAdminDBService{
		serverSettings:   serverSettings,
		profileDBService: profileDBService,
		emailService:     emailService,
	}
}
