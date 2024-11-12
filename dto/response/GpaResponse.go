package response

type GpaResponse struct {
	CurrentGpa   float64                  `json:"current_gpa"`
	PredictedGpa float64                  `json:"predicted_next_semester_gpa"`
	Semesters    []SemesterResultResponse `json:"semesters"`
}

type SemesterResultResponse struct {
	SemesterNumber       string                 `json:"semester_number"`
	AcademicYear         string                 `json:"academic_year"`
	Gpa                  float64                `json:"gpa"`
	TotalNumberOfCredits int                    `json:"total_number_of_credits"`
	CourseResults        []CourseResultResponse `json:"course_results"`
}

type CourseResultResponse struct {
	Name            string                   `json:"name"`
	NumberOfCredits int                      `json:"number_of_credits"`
	Gpa4Scale       float64                  `json:"gpa_4_scale"`
	Gpa10Scale      float64                  `json:"gpa_10_scale"`
	ComponentScore  []ComponentScoreResponse `json:"component_score"`
}

type ComponentScoreResponse struct {
	Name        string  `json:"name"`
	ScoreWeight float64 `json:"score_weight"`
	Score       float64 `json:"score"`
}
