package endpoint

import (
	"ScoreManagementSystem/dto/request"
	"ScoreManagementSystem/dto/response"
	"ScoreManagementSystem/model"
	"ScoreManagementSystem/service"
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
)

type StudentEndpoint interface {
	Login() endpoint.Endpoint
	Register() endpoint.Endpoint
	GetStudentInfo() endpoint.Endpoint
}

type studentEndpoint struct {
	studentService service.StudentService
}

func (s *studentEndpoint) Login() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		loginReq := req.(request.LoginRequest)
		validate := validator.New()
		err := validate.Struct(loginReq)
		if err != nil {
			return nil, err
		}

		var token string
		token, err = s.studentService.Login(ctx, loginReq)
		if err != nil {
			return nil, err
		}

		return response.LoginResponse{
			Token: token,
		}, nil
	}
}

func (s *studentEndpoint) Register() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		newStudent := req.(model.Student)
		validate := validator.New()
		err := validate.Struct(newStudent)
		if err != nil {
			return nil, err
		}
		err = s.studentService.Register(ctx, newStudent)
		if err != nil {
			return nil, err
		}
		return response.Message{Message: "Ok"}, nil
	}
}

func (s *studentEndpoint) GetStudentInfo() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		studentId := req.(string)
		student, err := s.studentService.GetStudentById(ctx, studentId)
		if err != nil {
			return nil, err
		}
		return response.GetStudentResponse{
			Id:             student.Id,
			Name:           student.Name,
			DateOfBirth:    student.DateOfBirth,
			Gender:         student.Gender,
			Email:          student.Email,
			IdentityNumber: student.IdentityNumber,
			PhoneNumber:    student.PhoneNumber,
			Address:        student.Address,
			Class:          student.Class,
			SchoolYear:     student.SchoolYear,
			FieldOfStudy:   student.FieldOfStudy,
		}, nil

	}
}

func NewStudentEndpoint(studentService service.StudentService) StudentEndpoint {
	return &studentEndpoint{
		studentService: studentService,
	}
}
