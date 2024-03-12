package authentication_db

import (
	"errors"
	"github.com/google/uuid"
	aws_ses "resume-review-api/aws-ses"
	email_templates "resume-review-api/email-templates"
	"resume-review-api/mongodb"
	"time"
)

func (s *ResumeAIAuthDBService) CreateForgotPassword(emailAddress string, ipAddress string, userAgent string) error {

	if emailAddress == "" {
		return errors.New("invalid email")
	}

	userId, err := s.profileDBService.GetUserIdByEmail(emailAddress)
	if err != nil {
		return err
	}

	// Check if Active User
	profile, err := s.profileDBService.GetProfileByUserId(userId)
	if err != nil {
		return err
	}

	if profile.ActiveUser == false {
		return errors.New("inactive user")
	}

	// Create Session DB
	token := uuid.New().String()
	doc := mongodb.ForgotPassword{
		Token:        token,
		UserId:       userId,
		CreationIP:   ipAddress,
		CreationDate: time.Now().UTC(),
		Expiration:   time.Now().UTC().Add(time.Hour * 24),
		Active:       true,
	}

	if _, err = mongodb.NewDocument(s.serverSettings.DBName, "forgot_passwords", doc); err != nil {
		return err
	}

	// Send Forgot Password Email
	s.emailService.SendEmailSES(
		email_templates.ForgotPasswordTemplate("https://"+s.serverSettings.BaseURL+"/resetpassword/"+token, userAgent, ipAddress),
		"Resume Reviewer - Forgot Password",
		s.serverSettings.FromEmail,
		aws_ses.Recipient{
			ToEmails: []string{emailAddress},
		},
	)

	// Done

	return nil
}
