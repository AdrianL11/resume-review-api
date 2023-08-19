package resume_routes

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
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
	Condense       bool   `json:"condense"`
	SkillCount     int    `json:"skill_count"`
	JobDutyCount   int    `json:"job_duty_count"`
}

func ReviewResume(c echo.Context) error {

	// Create Resume Review Bind
	var resumeReviewBind ResumeReviewBind
	if err := c.Bind(&resumeReviewBind); err != nil {
		log.Println("[Review Resume] Binding - " + err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate Binding
	if err := c.Validate(resumeReviewBind); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		log.Println("[Review Resume] Valid Session - " + err.Error())
		return c.NoContent(http.StatusUnauthorized)
	}

	// Create Resume
	mimeType, err := resume2.GetMimeType(resumeReviewBind.Resume)
	if err != nil {
		log.Println("[Review Resume] Create Resume - " + err.Error())
		return err
	}

	resume, err := resume2.ConvertToPlainText(resumeReviewBind.Resume, mimeType)
	if err != nil {
		log.Println("[Review Resume] Convert Resume - " + err.Error())
		return err
	}

	// Create Job Description
	mimeType, err = resume2.GetMimeType(resumeReviewBind.JobDescription)
	if err != nil {
		log.Println("[Review Resume] Create JD - " + err.Error())
		return err
	}

	jobDescription, err := resume2.ConvertToPlainText(resumeReviewBind.JobDescription, mimeType)
	if err != nil {
		log.Println("[Review Resume] Convert JD - " + err.Error())
		return err
	}

	// Get User ID
	userId, err := mongodb.GetProfileBySession(c)
	if err != nil {
		log.Println("[Review Resume] Get User ID - " + err.Error())
		return err
	}

	// Start Response Time
	startTime := time.Now().UTC()

	// Get JSON Object
	var jsonObject = "{  \"name\": \"\",  \"email\": \"\",  \"phone\": \"\",  \"address\": \"\",  \"summary\": \"\",  \"skills\": [],  \"experiences\": [    {      \"company\": \"\",      \"location\": \"\",      \"dates\": \"\",      \"job_title\": \"\",      \"duties\": []    }  ],  \"educations\": [    {      \"school_name\": \"\",      \"location\": \"\",      \"type\": \"\",      \"graduation\": \"\"    }  ],  \"score\": \"\",  “scoreReason: \"\", \"recruiter_summary\": \"\"}"
	var tasks = fmt.Sprintf("1. CORRECT grammar, punctuation, and spelling mistakes. \n2. Correct capitalization mistakes.\n3. Format dates with full spelling of the month and correct and SPACING consistent.\n4. Return ALL education certifications. List degrees from highest to lowest. Return maximum 10 most recent work history in chronological order regardless of job description. \n5. List a maximum of %d job duties per work history.\n6. List a maximum of %d skills. If skills are empty, generate %d skills that align with the job description.\n7. Give a final scoring of resume against the job description given and any suggestions that may make the candidate look better.\n8. Create a detailed summary with experience and key accomplishments based on resume and job description and use gender-neutral terms that a recruiter can send to a hiring manager.\n9. Always return a recruiter summary.\n10. Put into given json object.\n11. limit prose.\n12. only provide json, no explanation.\n", resumeReviewBind.JobDutyCount, resumeReviewBind.SkillCount, resumeReviewBind.SkillCount)
	var options = ""

	if resumeReviewBind.Condense {
		options = "1. Condense Resume\n"
	} else {
		options = "1. DO NOT CONDENSE Resume\n"
	}

	// Tests
	fmt.Println("Job Description: \n\n" + jobDescription)
	fmt.Println("Resume: \n\n" + resume)
	fmt.Println("JSON Object: \n\n" + jsonObject)
	fmt.Println("Tasks: \n\n" + tasks)
	fmt.Println("Options: \n\n" + options)

	// Set OpenAI Prompts
	ret, err := resume2.CreateGPTRequest([]resume2.Message{
		{
			Role:    "system",
			Content: "You are a very experienced resume reviewer with over 20 years of experience. Read carefully and be very thorough. This is a important task. Be professional. Ensure that you return all work history.",
		},
		{
			Role:    "user",
			Content: "Job Description: \n\n" + jobDescription,
		},
		{
			Role:    "user",
			Content: "Resume: \n\n" + resume,
		},
		{
			Role:    "user",
			Content: "JSON Object: \n\n" + jsonObject,
		},
		{
			Role:    "user",
			Content: "Tasks: \n\n" + tasks,
		},
		{
			Role:    "user",
			Content: "Tasks: \n\n" + options,
		},
		{
			Role:    "assistant",
			Content: "{}",
		},
		{
			Role:    "user",
			Content: "You didnt do all the work experience. Please make sure to include every work experience. Send only the JSON object.",
		},
	})

	fmt.Println(ret)

	if err != nil {
		log.Println("[Review Resume] CreateGPTRequest - " + err.Error())
		return err
	}

	// Get End Time & Final Response Time
	endTime := time.Now().UTC()
	responseTime := endTime.Sub(startTime).Seconds()

	// Create Return
	var jsonObj resume2.JSONObject
	err = json.Unmarshal([]byte(ret), &jsonObj)
	if err != nil {
		log.Println("[Review Resume] Create Return - " + err.Error())
		return err
	}

	// Insert Into Database
	err = resume_db.InsertResumeReview(userId.ID, jsonObj, responseTime)
	if err != nil {
		log.Println("[Review Resume] Insert into DB - " + err.Error())
		return err
	}

	// Done
	return c.JSON(http.StatusOK, jsonObj)
}
