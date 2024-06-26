package authentication_db

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"os"
	aws_ses "resume-review-api/aws-ses"
	email_templates "resume-review-api/email-templates"
	"resume-review-api/mongodb"
	"time"
)

func CreateForgotPassword(c echo.Context, emailAddress string) error {

	if emailAddress == "" {
		return errors.New("invalid email")
	}

	userId, err := mongodb.GetUserIdByEmail(emailAddress)
	if err != nil {
		return err
	}

	// Check if Active User
	profile, err := mongodb.GetProfilebyUserId(userId)
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
		CreationIP:   c.RealIP(),
		CreationDate: time.Now().UTC(),
		Expiration:   time.Now().UTC().Add(time.Hour * 24),
		Active:       true,
	}

	if _, err = mongodb.NewDocument(os.Getenv("db_name"), "forgot_passwords", doc); err != nil {
		return err
	}

	// Send Forgot Password Email
	aws_ses.SendEmailSES(
		email_templates.ForgotPasswordTemplate("https://"+os.Getenv("base_url")+"/resetpassword/"+token, c.Request().UserAgent(), c.RealIP()),
		"Resume Reviewer - Forgot Password",
		os.Getenv("from_email"),
		aws_ses.Recipient{
			ToEmails: []string{emailAddress},
		},
	)

	// Done

	return nil
}
