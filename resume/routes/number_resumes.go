package resume_routes

import (
	"net/http"
	"resume-review-api/mongodb"
	"resume-review-api/resume"
	resume_db "resume-review-api/resume/database"
	session_db "resume-review-api/session/database"

	"github.com/labstack/echo/v4"
)

func GetResumeMatched(c echo.Context) error {

	// Is Session Valid
	err := session_db.ValidateSession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// Get User ID
	profile, err := mongodb.GetProfileBySession(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	resumeCount, jobMatch, err := resume_db.GetResumeMatched(profile.ID)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, resume.ReturnJobMatched{
		ResumeCount: resumeCount,
		JobMatch:    jobMatch,
	})
}
