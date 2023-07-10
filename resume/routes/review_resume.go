package resume_routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/dslipak/pdf"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	resume2 "resume-review-api/resume"
	session_db "resume-review-api/session/database"
	"strings"
)

type ResumeReviewBind struct {
	Resume         string `json:"resume"`
	JobDescription string `json:"job_description"`
}

func ReviewResume(c echo.Context) error {

	// Create Resume Review Bind
	var resumeReviewBind ResumeReviewBind
	if err := c.Bind(&resumeReviewBind); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(resumeReviewBind); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Create Resume
	resume, err := ParseResume(resumeReviewBind.Resume)
	if err != nil {
		return err
	}

	// Get JSON Object
	var jsonObject = "{  \"name\": \"\",  \"email\": \"\",  \"phone\": \"\",  \"address\": \"\",  \"summary\": \"\",  \"skills\": [],  \"experiences\": [    {      \"company\": \"\",      \"location\": \"\",      \"dates\": \"\",      \"job_title\": \"\",      \"duties\": []    }  ],  \"educations\": [    {      \"school_name\": \"\",      \"location\": \"\",      \"type\": \"\",      \"graduation\": \"\"    }  ],  \"score\": \"\",  “scoreReason: \"\", \"recruiter_summary\": \"\"}"
	var prompt = "redo resume in give json object with corrected grammar, punctuation, corrected formatting, and with only 4-5 duties per experience, also list all educations and certifications. Make sure all dates are formatted the same. If skills are empty, generate 6 skills, otherwise give me the top 6 based on job description.  Give a final scoring of resume against the job description given. also create a short paragraph summary based on resume and job description. limit prose. only provide json, no explanation."
	var builtPrompt = jsonObject + "\n\nResume: \n" + resume + "\n\nJob Description: \n" + resumeReviewBind.JobDescription + "\n\n\n" + prompt

	// Set OpenAI Prompts
	ret, err := resume2.CreateGPTRequest([]resume2.Message{
		{
			Role:    "user",
			Content: builtPrompt,
		},
	})
	if err != nil {
		return err
	}

	// Create Return
	var jsonObj resume2.JSONObject
	err = json.Unmarshal([]byte(ret), &jsonObj)
	if err != nil {
		return err
	}

	// Done
	return c.JSON(http.StatusOK, jsonObj)
}

func ParseResume(res string) (string, error) {

	if !strings.Contains(res, "data:application/pdf;base64,") {
		return "", errors.New("no pdf found")
	}

	var resume = strings.Replace(res, "data:application/pdf;base64,", "", -1)
	var _uuid = uuid.New().String()

	dec, err := base64.StdEncoding.DecodeString(resume)
	if err != nil {
		return "", err
	}

	f, err := os.Create(_uuid + ".pdf")
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	if _, err := f.Write(dec); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}

	// Read PDF
	r, err := pdf.Open(_uuid + ".pdf")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)

	err = os.Remove(_uuid + ".pdf")
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
