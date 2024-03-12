package aws_ses

import "resume-review-api/util/resume_ai_env"

type EmailService struct {
	serverSettings resume_ai_env.ServerSettings
}

func NewEmailService(serverSettings resume_ai_env.ServerSettings) *EmailService {
	return &EmailService{
		serverSettings: serverSettings,
	}
}
