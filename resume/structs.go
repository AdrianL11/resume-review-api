package resume

type JSONObject struct {
	Name             string       `json:"name"`
	Email            string       `json:"email"`
	Phone            string       `json:"phone"`
	Address          string       `json:"address"`
	Summary          string       `json:"summary"`
	Skills           []string     `json:"skills"`
	Experiences      []Experience `json:"experiences"`
	Educations       []Education  `json:"educations"`
	Score            string       `json:"score"`
	ScoreREason      string       `json:"scoreReason"`
	RecruiterSummary string       `json:"recruiter_summary"`
}

type Experience struct {
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Dates    string   `json:"dates"`
	JobTitle string   `json:"job_title"`
	Duties   []string `json:"duties"`
}

type Education struct {
	SchoolName     string `json:"school_name"`
	Location       string `json:"location"`
	Type           string `json:"type"`
	GraduationYear string `json:"graduation"`
}
