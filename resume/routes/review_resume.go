package resume_routes

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"resume-review-api/mongodb"
	resume2 "resume-review-api/resume"
	resume_db "resume-review-api/resume/database"
	session_db "resume-review-api/session/database"
	"time"
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
	resume, err := resume2.ParseResume(resumeReviewBind.Resume)
	if err != nil {
		return err
	}

	// Get User ID
	userId, err := mongodb.GetProfileBySession(c)
	if err != nil {
		return err
	}

	// Start Response Time
	startTime := time.Now().UTC()

	// Get JSON Object
	var jsonObject = "{  \"name\": \"\",  \"email\": \"\",  \"phone\": \"\",  \"address\": \"\",  \"summary\": \"\",  \"skills\": [],  \"experiences\": [    {      \"company\": \"\",      \"location\": \"\",      \"dates\": \"\",      \"job_title\": \"\",      \"duties\": []    }  ],  \"educations\": [    {      \"school_name\": \"\",      \"location\": \"\",      \"type\": \"\",      \"graduation\": \"\"    }  ],  \"score\": \"\",  “scoreReason: \"\", \"recruiter_summary\": \"\"}"
	var prompt = "redo resume in give json object with corrected grammar, punctuation, corrected formatting, and with only 4-5 duties per experience, also list all educations and certifications. Make sure all dates are formatted the same. If skills are empty, generate 6 skills, otherwise give me the top 6 based on job description.  Give a final scoring of resume against the job description given and any suggestions that may make the candidate look better. also create a detailed summary with experience and key accomplishments based on resume and job description. limit prose. only provide json, no explanation."
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

	// Get End Time & Final Response Time
	endTime := time.Now().UTC()
	responseTime := endTime.Sub(startTime).Seconds()

	// Create Return
	var jsonObj resume2.JSONObject
	err = json.Unmarshal([]byte(ret), &jsonObj)
	if err != nil {
		return err
	}

	// Insert Into Database
	err = resume_db.InsertResumeReview(userId.ID, jsonObj, responseTime)
	if err != nil {
		return err
	}

	// Done
	return c.JSON(http.StatusOK, jsonObj)
}
