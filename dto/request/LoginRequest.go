package request

type LoginRequest struct {
	StudentId string `json:"student_id" validate:"required"`
	Password  string `json:"password" validate:"required"`
}
