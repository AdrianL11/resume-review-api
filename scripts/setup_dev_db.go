package scripts

import (
	"errors"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"resume-review-api/mongodb"
	"resume-review-api/util"
	"resume-review-api/util/resume_ai_env"
	"time"
)

const DefaultOwnerEmail = "owner@vdart.ai"

func SetupDevEnvIfNeeded() error {
	if !resume_ai_env.IsDev() {
		return errors.New("expected development environment")
	}

	// Admin already exists?
	existingAdmin, err := mongodb.GetUserIdByEmail(DefaultOwnerEmail)
	if !existingAdmin.IsZero() {
		return nil
	}

	password, err := util.HashPassword("password")
	if err != nil {
		return err
	}

	defaultUser := bson.D{
		{"email_address", DefaultOwnerEmail},
		{"role", mongodb.OwnerRole},
		{"creation_date", time.Now().UTC()},
		{"first_login", true},
		{"active_user", true},
		{"password", password},
	}
	_, err = mongodb.NewDocument(resume_ai_env.GetSettingsForEnv().DBName, "users", defaultUser)
	if err != nil {
		return err
	}
	log.Info("Created default user %s", DefaultOwnerEmail)

	return nil
}
