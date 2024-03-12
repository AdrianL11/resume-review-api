package resume_ai_env

import (
	"log"
	"os"
	"sync"
)

type ServerSettings struct {
	APIVersion          string
	BaseURL             string
	FromEmail           string
	GPTVersion          string
	DBName              string
	DBURL               string
	DBUsername          string
	DBPassword          string
	OpenAIKey           string
	SESAkid             string
	SESSecret           string
	SessionKey          string
	SessionCookieName   string
	SessionCookieDomain string
	SpacesKey           string
	SpacesSecret        string
	AzureDeploymentName string
	AzureKey            string
	Port                string
}

var settings ServerSettings
var once sync.Once

func GetSettingsForEnv() ServerSettings {
	once.Do(func() {
		settings = ServerSettings{
			APIVersion:          getEnvOrPanic("api_version"),
			BaseURL:             getEnvOrPanic("base_url"),
			FromEmail:           getEnvOrPanic("from_email"),
			GPTVersion:          getEnvOrPanic("gpt_version"),
			DBName:              getEnvOrPanic("db_name"),
			DBURL:               getEnvOrPanic("mongodb_url"),
			DBUsername:          getEnvOrPanic("mongodb_username"),
			DBPassword:          getEnvOrPanic("mongodb_password"),
			OpenAIKey:           getEnvOrPanic("openai_key"),
			SESAkid:             getEnvOrPanic("ses_akid"),
			SESSecret:           getEnvOrPanic("ses_secret"),
			SessionKey:          getEnvOrPanic("session_key"),
			SessionCookieName:   getEnvOrPanic("session_name"),
			SessionCookieDomain: getEnvOrPanic("session_cookie_domain"),
			SpacesKey:           getEnvOrPanic("spaces_key"),
			SpacesSecret:        getEnvOrPanic("spaces_secret"),
			AzureDeploymentName: getEnvOrPanic("azure_deployment_name"),
			AzureKey:            getEnvOrPanic("azure_key"),
			Port:                getEnvOrPanic("PORT"),
		}
	})
	return settings
}

func getEnvOrPanic(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatal("MISSING ENV: %s", name)
	}
	return val
}
