package model

type Student struct {
	Id             string `json:"id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	DateOfBirth    string `json:"date_of_birth" validate:"required"`
	Gender         string `json:"gender" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	IdentityNumber string `json:"identity_number" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required"`
	Address        string `json:"address" validate:"required"`
	Password       string `json:"password" validate:"required"`
	Class          string `json:"class" validate:"required"`
	SchoolYear     string `json:"school_year" validate:"required"`
	FieldOfStudy   string `json:"field_of_study" validate:"required"`
}
