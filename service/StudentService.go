package service

import (
	"ScoreManagementSystem/dto/request"
	"ScoreManagementSystem/model"
	"ScoreManagementSystem/repo"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type StudentService interface {
	Login(ctx context.Context, req request.LoginRequest) (string, error)
	Register(ctx context.Context, req model.Student) error
	GetStudentById(ctx context.Context, studentId string) (model.Student, error)
}

type studentService struct {
	studentRepo      repo.StudentRepo
	jwtService       JwtService
	redisStudentRepo repo.RedisStudentRepo
}

func (s *studentService) Login(ctx context.Context, req request.LoginRequest) (string, error) {
	student, err := s.GetStudentById(ctx, req.StudentId)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("incorrect password")
	}
	return s.jwtService.CreateToken(student.Id)
}
func (s *studentService) Register(ctx context.Context, req model.Student) error {
	var hash []byte
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return errors.New("failed to hash password")
	}
	req.Password = string(hash)

	return s.studentRepo.AddNewStudent(ctx, req)
}
func (s *studentService) GetStudentById(ctx context.Context, studentId string) (model.Student, error) {
	//Get student info in redis if exist
	res, err := s.redisStudentRepo.GetStudentInfoByStudentId(ctx, studentId)
	if err == nil {
		return res, nil
	}
	//Get student info from database
	student, err := s.studentRepo.GetStudentById(ctx, studentId)
	if err != nil {
		return model.Student{}, err
	}

	//Put student info in redis
	s.redisStudentRepo.SaveStudentInfo(ctx, student)
	return student, nil
}

func NewStudentService(studentRepo repo.StudentRepo, jwtService JwtService, redisStudentRepo repo.RedisStudentRepo) StudentService {
	return &studentService{
		studentRepo:      studentRepo,
		jwtService:       jwtService,
		redisStudentRepo: redisStudentRepo,
	}
}
