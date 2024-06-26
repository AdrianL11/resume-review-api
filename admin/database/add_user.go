package admin_db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	aws_ses "resume-review-api/aws-ses"
	email_templates "resume-review-api/email-templates"
	"resume-review-api/mongodb"
	"time"
)

func AddUser(currentUser primitive.ObjectID, emailAddress string, role mongodb.Role) error {

	// Does Email Address Already Exist?
	profile, err := mongodb.GetUserIdByEmail(emailAddress)
	if err == nil && profile.String() != "" {
		return errors.New("user already exists")
	}

	// Profile Doesn't Exist, Add
	doc := bson.D{
		{"email_address", emailAddress},
		{"role", role},
		{"created_by", currentUser},
		{"creation_date", time.Now().UTC()},
		{"first_login", true},
		{"active_user", true},
	}

	result, err := mongodb.NewDocument(os.Getenv("db_name"), "users", doc)
	if err != nil {
		return err
	}

	// Send Email
	objId := result.InsertedID.(primitive.ObjectID)

	aws_ses.SendEmailSES(email_templates.NewUserEmail("https://"+os.Getenv("base_url")+"/acceptinvite/"+objId.Hex()), "You have been invited to Resume Reviewer!", "no-reply@vdart.ai", aws_ses.Recipient{
		ToEmails: []string{emailAddress},
	})

	return nil
}
