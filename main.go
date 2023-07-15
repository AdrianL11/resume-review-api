package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	admin_routes "resume-review-api/admin/routes"
	authentication_routes "resume-review-api/authentication/routes"
	profile_routes "resume-review-api/profile/routes"
	resume_routes "resume-review-api/resume/routes"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	// Create New Echo Server
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Create Session Store
	e.Use(session.MiddlewareWithConfig(session.Config{
		Store: sessions.NewCookieStore([]byte(os.Getenv("session_key"))),
	}))

	// CORS Setup
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodPost, http.MethodGet},
		AllowOrigins: []string{"http://" + os.Getenv("base_url"), "https://" + os.Getenv("base_url")},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderAccessControlAllowOrigin,
		},
		AllowCredentials: true,
	}))

	// Base
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Version: %s", os.Getenv("api_version")))
	})

	// Authentication Routes
	e.POST("/login", authentication_routes.Login)
	e.GET("/logout", authentication_routes.Logout)
	e.POST("/forgot_password", authentication_routes.ForgotPassword)
	e.POST("/forgot_password_validate", authentication_routes.ForgotPasswordValidate)
	e.POST("/reset_password", authentication_routes.ResetPassword)
	e.GET("/logged_in", authentication_routes.LoggedIn)

	// Profile Routes
	e.GET("/profile", profile_routes.GetProfile)
	e.POST("/profile/set", profile_routes.SetProfile)
	e.GET("/profile/sessions", profile_routes.GetActiveSessions)
	e.POST("/profile/update", profile_routes.UpdateProfile)
	e.POST("/profile/change_password", profile_routes.ChangePassword)
	e.POST("/new_user_validate", profile_routes.NewUserValidate)

	// Admin Routes
	e.GET("/admin/get_profiles", admin_routes.GetUsers)
	e.POST("/admin/get_user_sessions", admin_routes.GetUserSessions)
	e.POST("/admin/get_profile", admin_routes.GetProfileById)
	e.POST("/admin/update_profile", admin_routes.UpdateProfile)
	e.POST("/admin/deactivate_user", admin_routes.DeactivateUser)
	e.POST("/admin/add_user", admin_routes.AddUser)

	// Resume Routes
	e.POST("/resume/review", resume_routes.ReviewResume)
	e.GET("/resume/counts", resume_routes.GetResumeCountInfo)

	// Set Server Port
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8085"
	}

	// Start Server
	err := e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Error Starting Server: " + err.Error())
	}

	/*
		Paths:

			Session Cookie Data: {session_id, object_id}

			POST /login
				Params: 	 {email_address, password}
				Description: User POST email_address and password. API checks DB('resume_reviewer') TABLE('users') for
							 email_address and password. If found, create a session cookie and send to user. Insert
							 session cookie into DB('resume_reviewer') TABLE('sessions'). Return JSON 'success'. If not
							 found, return JSON, 'User not found'.

			POST /forgot_password
				Params: 	 {email_address}
				Description: User POST email_address. API checks DB('resume_reviewer') TABLE('users') for email_address.
							 If found, API inserts into DB('resume_reviewer') TABLE('forgot_passwords'). The data into
							 table is {uuid, object_id, creation_date, creation_ip, expiration_date, active}.
							 API then emails the email_address a forgot password link. EX: /forgot_password?o={uuid}

			POST /forgot_password_check
				Params: 	 {uuid}
				Description: User POST uuid. API checks if uuid exists in DB('resume_reviewer')
							 TABLE('forgot_passwords'). If uuid exists, API checks if expired. If not expired and
							 exists, return JSON 'success'. If expired or does not exist, return JSON 'error'.

			POST /reset_password
				Params:		 {uuid, new_password}
				Description: User POST uuid and new_password. API checks if uuid exists in DB('resume_reviewer')
							 TABLE('forgot_passwords'). If uuid exists, API checks if expired. If not expired and
							 exists, API gets object_id from document and inserts into DB('resume_reviewer')
							 TABLE('users'), new_password. API then inserts into DB('resume_reviewer')
							 TABLE('forgot_passwords'), active = false.

			POST /set_profile
				Params: 	 {object_id, password, first_name, last_name, country, profile_image}
				Description: User POST object_id, password, first_name, last_name, country. API checks DB('resume_reviewer')
							 TABLE('users') if object_id exists. If object_id exists and first_login = true, API inserts
							 into DB('resume_reviewer') TABLE('users') password = password, first_login = false, first_name,
							 last_name, country, profile_image.

			POST /new_user_check
				Params:		 {object_id}
				Description: User POST object_id. API checks DB('resume_reviewer') TABLE('users') if object_id exists.
							 If object_id exists and first_login = true, API return 'true', else returns 'false'.
							 This will check if user can set new profile and password on this object_id.

			GET /profile
				Params:		 {}
				Description: User GET. API checks session cookie. If cookie exists, API checks DB('resume_reviewer')
							 TABLE('sessions'). If session exists and is active, API gets profile_id from session cookie
							 and returns {first_name, last_name, country, role}. If session or cookie does not exist,
							 return JSON 'error'.

			POST /change_password
				Params:		 {old_password, new_password}
				Description: User POST {old_password, new_password}. API checks current session. If session is valid from
							 cookie and MongoDB, then API gets profile_id from session. With profile_id, MOngoDB gets
							 user, which is the object_id for the user database. API checks old_password with DB password.
							 If old_password matches, API edits and inserts new_password. Returns JSON 'success'. Otherwise,
							 JSON return 'error'.

			POST /add_user
				Params: 	 {email_address, role}
				Description: User POST {email_address, role}. API checks session cookie. If valid, continue. If invalid,
							 return JSON 'error'. API checks if current user role >= new user role. If true, API will
							 insert into DB('resume_reviewer') TABLE('users') {email_address, first_login=true, uuid, role,
							 active = true}. If session is invalid, or role < admin, or the current user role is not >=
							 new user role, return JSON 'error'. API returns JSON 'uuid' so that frontend can send email.
							 API also checks if email_address exists, if yes, return 'error'.

			POST /deactivate_user
				Params: 	 {email_address}
				Description: User POST {email address}. API checks session cookie. If valid, continue. If invalid,
							 return JSON 'error'. API checks if current user role >= deactivated user role. and current
							 user is admin or greater. If true, API will edit DB('resume_reviewer') TABLE('users')
							 {active = false}. Return JSON 'success'. If {active} is already false, return JSON 'error'.

			GET /get_users

			POST /review_resume
				Params: {resume, job_description}

			GET /logout

	*/
}
