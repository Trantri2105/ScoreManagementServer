package response

type GetStudentResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	DateOfBirth    string `json:"date_of_birth"`
	Gender         string `json:"gender"`
	Email          string `json:"email"`
	IdentityNumber string `json:"identity_number"`
	PhoneNumber    string `json:"phone_number"`
	Address        string `json:"address"`
	Class          string `json:"class"`
	SchoolYear     string `json:"school_year"`
	FieldOfStudy   string `json:"field_of_study"`
}
