package profile_util

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"resume-review-api/util/resume_ai_env"
	"strings"
)

func GetImageCDNURL(base64Image string) (string, error) {

	// Check if Image is blank?
	if base64Image == "" {
		return "", nil
	}

	// Step 1: Get Image Info
	imageType := strings.Split(base64Image, ";base64,")[0]
	image, _ := base64.StdEncoding.DecodeString(strings.Split(base64Image, ";base64,")[1])
	mimeType := ""

	if imageType == "data:image/png" {
		mimeType = "png"
	} else if imageType == "data:image/jpeg" {
		mimeType = "jpg"
	}

	imageLink := uuid.New().String() + "." + mimeType

	// Step 2: Define the parameters for the session you want to create.
	key := resume_ai_env.GetSettingsForEnv().SpacesKey       // Access key pair. You can create access key pairs using the control panel or API.
	secret := resume_ai_env.GetSettingsForEnv().SpacesSecret // Secret access key defined through an environment variable.

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""), // Specifies your credentials.
		Endpoint:         aws.String("https://nyc3.digitaloceanspaces.com"), // Find your endpoint in the control panel, under Settings. Prepend "https://".
		S3ForcePathStyle: aws.Bool(false),                                   // // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
		Region:           aws.String("nyc3"),                                // Must be "us-east-1" when creating new Spaces. Otherwise, use the region in your endpoint, such as "nyc3".

	}

	// Step 3: The new session validates your request and directs it to your Space's specified endpoint using the AWS SDK.
	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	// Step 4: Define the parameters of the object you want to upload.
	object := s3.PutObjectInput{
		Bucket: aws.String("resume-reviewer-cdn"),         // The path to the directory you want to upload the object to, starting with your Space name.
		Key:    aws.String("profile-images/" + imageLink), // Object key, referenced whenever you want to access this file later.
		Body:   strings.NewReader(string(image)),          // The object's contents.
		ACL:    aws.String("public-read"),                 // Defines Access-control List (ACL) permissions, such as private or public.

	}

	// Step 5: Run the PutObject function with your parameters, catching for errors.
	_, err := s3Client.PutObject(&object)
	if err != nil {
		return "", err
	}

	return "https://resume-reviewer-cdn.nyc3.cdn.digitaloceanspaces.com/profile-images/" + imageLink, nil
}
