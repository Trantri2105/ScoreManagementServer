package repo

import (
	"ScoreManagementSystem/model"
	"context"
	"database/sql"
	"github.com/lib/pq"
	"log"
)

type GpaRepo interface {
	GetCourseResultByStudentId(ctx context.Context, studentId string) ([]model.CourseResult, error)
	GetScoreComponentScore(ctx context.Context, courseIds []int) ([]model.ComponentScore, error)
}

type gpaRepo struct {
	db *sql.DB
}

func (g *gpaRepo) GetCourseResultByStudentId(ctx context.Context, studentId string) ([]model.CourseResult, error) {
	query := `SELECT c.name, c.number_of_credit,cr.id, cr.semester_number, cr.academic_year
			  FROM courses as c 
			  JOIN course_results as cr ON cr.course_id=c.id
			  WHERE cr.student_id = $1`

	rows, err := g.db.QueryContext(ctx, query, studentId)
	if err != nil {
		log.Println("Gpa repo, get course result error :", err)
		return nil, err
	}

	courseResults := make([]model.CourseResult, 0)
	for rows.Next() {
		var courseResult model.CourseResult
		err = rows.Scan(&courseResult.Name, &courseResult.NumberOfCredits, &courseResult.Id, &courseResult.SemesterNumber, &courseResult.AcademicYear)
		if err != nil {
			log.Println("Gpa repo, get course result error :", err)
			return nil, err
		}
		courseResults = append(courseResults, courseResult)
	}
	return courseResults, nil
}

func (g *gpaRepo) GetScoreComponentScore(ctx context.Context, courseIds []int) ([]model.ComponentScore, error) {
	query := `SELECT name, score_weight, score, course_result_id
			  FROM component_scores
			  WHERE course_result_id = ANY($1)`

	rows, err := g.db.QueryContext(ctx, query, pq.Array(courseIds))
	if err != nil {
		log.Println("Gpa repo, get score component score error :", err)
		return nil, err
	}
	scoreComponents := make([]model.ComponentScore, 0)
	for rows.Next() {
		var scoreComponent model.ComponentScore
		err = rows.Scan(&scoreComponent.Name, &scoreComponent.ScoreWeight, &scoreComponent.Score, &scoreComponent.CourseResultId)
		if err != nil {
			log.Println("Gpa repo, get score component score error :", err)
			return nil, err
		}
		scoreComponents = append(scoreComponents, scoreComponent)
	}
	return scoreComponents, nil
}

func NewGpaRepo(db *sql.DB) GpaRepo {
	return &gpaRepo{db: db}
}
