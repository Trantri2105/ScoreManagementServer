package repo

import (
	"ScoreManagementSystem/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type StudentRepo interface {
	GetStudentById(ctx context.Context, id string) (model.Student, error)
	AddNewStudent(ctx context.Context, student model.Student) error
}

type studentRepo struct {
	db *sql.DB
}

const StudentTable = "students"

func (s *studentRepo) GetStudentById(ctx context.Context, id string) (model.Student, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, StudentTable)
	row, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		log.Println("Student repo, get student by id err: ", err)
		return model.Student{}, err
	}
	student := model.Student{}
	ok := row.Next()
	if !ok {
		return model.Student{}, errors.New("student not found")
	}
	err = row.Scan(&student.Id,
		&student.Name,
		&student.DateOfBirth,
		&student.Gender,
		&student.Email,
		&student.IdentityNumber,
		&student.PhoneNumber,
		&student.Address,
		&student.Password,
		&student.Class,
		&student.SchoolYear,
		&student.FieldOfStudy)
	if err != nil {
		log.Println("Student repo, get student by id err: ", err)
		return model.Student{}, err
	}
	return student, nil
}

func (s *studentRepo) AddNewStudent(ctx context.Context, student model.Student) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, name, date_of_birth, gender, email, identity_number,
              phone_number, address, password, class, school_year, field_of_study)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, StudentTable)
	_, err := s.db.ExecContext(ctx, query,
		student.Id,
		student.Name,
		student.DateOfBirth,
		student.Gender,
		student.Email,
		student.IdentityNumber,
		student.PhoneNumber,
		student.Address,
		student.Password,
		student.Class,
		student.SchoolYear,
		student.FieldOfStudy)
	return err
}

func NewStudentRepo(db *sql.DB) StudentRepo {
	return &studentRepo{db: db}
}
