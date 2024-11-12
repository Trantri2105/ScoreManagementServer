package model

import "ScoreManagementSystem/dto/response"

type CourseResult struct {
	Id              int
	Name            string
	NumberOfCredits int
	SemesterNumber  string
	AcademicYear    string
	ComponentScores []response.ComponentScoreResponse
	Gpa10Scale      float64
	Gpa4Scale       float64
}
