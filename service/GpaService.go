package service

import (
	"ScoreManagementSystem/client"
	"ScoreManagementSystem/dto/request"
	"ScoreManagementSystem/dto/response"
	"ScoreManagementSystem/model"
	"ScoreManagementSystem/repo"
	"context"
	"fmt"
	"log"
	"sort"
)

type GpaService interface {
	GetStudentGpaByStudentId(ctx context.Context, studentId string) (response.GpaResponse, error)
	GetPredictedGpa(ctx context.Context, gpa []float64) float64
}

type gpaService struct {
	gpaRepo      repo.GpaRepo
	client       *client.ClientEndpoint
	redisGpaRepo repo.RedisGpaRepo
}

func (g *gpaService) GetStudentGpaByStudentId(ctx context.Context, studentId string) (response.GpaResponse, error) {
	//Get student gpa in redis if exist
	res, err := g.redisGpaRepo.GetStudentGpaByStudentId(ctx, studentId)
	if err == nil {
		return res, nil
	}
	//Get course result info
	courseResults, err := g.gpaRepo.GetCourseResultByStudentId(ctx, studentId)
	if err != nil {
		return response.GpaResponse{}, err
	}
	//Get component score info
	courseIds := make([]int, len(courseResults))
	for i, courseResult := range courseResults {
		courseIds[i] = courseResult.Id
	}
	componentScores, err := g.gpaRepo.GetScoreComponentScore(ctx, courseIds)
	if err != nil {
		return response.GpaResponse{}, err
	}
	//Calculate gpa for each course
	gpaResponse, err := g.calculateGpa(ctx, courseResults, componentScores)
	if err != nil {
		return response.GpaResponse{}, err
	}

	//Save gpa in cache
	g.redisGpaRepo.SaveStudentGpa(ctx, studentId, gpaResponse)
	return gpaResponse, nil

}

func convertTo4PointScale(grade10 float64) float64 {
	if grade10 >= 9.0 && grade10 <= 10.0 {
		return 4.0 // A+
	} else if grade10 >= 8.5 && grade10 < 9.0 {
		return 3.7 // A
	} else if grade10 >= 8.0 && grade10 < 8.5 {
		return 3.5 // B+
	} else if grade10 >= 7.0 && grade10 < 8.0 {
		return 3.0 // B
	} else if grade10 >= 6.5 && grade10 < 7.0 {
		return 2.5 // C+
	} else if grade10 >= 5.5 && grade10 < 6.5 {
		return 2.0 // C
	} else if grade10 >= 5.0 && grade10 < 5.5 {
		return 1.5 // D+
	} else if grade10 >= 4.0 && grade10 < 5.0 {
		return 1.0 // D
	} else if grade10 < 4.0 {
		return 0.0 // F
	} else {
		return -1
	}
}

func (g *gpaService) calculateGpa(ctx context.Context, courseResults []model.CourseResult, componentScores []model.ComponentScore) (response.GpaResponse, error) {
	//Calculate gpa for each course
	m := make(map[int]int)
	for i, courseResult := range courseResults {
		courseResult.ComponentScores = make([]response.ComponentScoreResponse, 0)
		m[courseResult.Id] = i
	}
	for _, componentScore := range componentScores {
		courseResults[m[componentScore.CourseResultId]].Gpa10Scale += componentScore.Score * componentScore.ScoreWeight
		res := response.ComponentScoreResponse{
			Name:        componentScore.Name,
			ScoreWeight: componentScore.ScoreWeight,
			Score:       componentScore.Score,
		}
		courseResults[m[componentScore.CourseResultId]].ComponentScores = append(courseResults[m[componentScore.CourseResultId]].ComponentScores, res)
	}
	for i := range courseResults {
		courseResults[i].Gpa4Scale = convertTo4PointScale(courseResults[i].Gpa10Scale)
	}

	//Calculate gpa for each semester
	semester := make([]response.SemesterResultResponse, 0)
	m2 := make(map[string]int)
	for _, courseResult := range courseResults {
		key := fmt.Sprintf("%s#%s", courseResult.SemesterNumber, courseResult.AcademicYear)
		_, ok := m2[key]
		if !ok {
			m2[key] = len(semester)
			semester = append(semester, response.SemesterResultResponse{
				SemesterNumber: courseResult.SemesterNumber,
				AcademicYear:   courseResult.AcademicYear,
			})
		}
		idx := m2[key]
		res := response.CourseResultResponse{
			Name:            courseResult.Name,
			NumberOfCredits: courseResult.NumberOfCredits,
			Gpa4Scale:       courseResult.Gpa4Scale,
			Gpa10Scale:      courseResult.Gpa10Scale,
			ComponentScore:  courseResult.ComponentScores,
		}
		semester[idx].CourseResults = append(semester[idx].CourseResults, res)
		semester[idx].Gpa += res.Gpa4Scale * float64(res.NumberOfCredits)
		semester[idx].TotalNumberOfCredits += res.NumberOfCredits

	}
	sort.Slice(semester, func(i, j int) bool {
		if semester[i].AcademicYear == semester[j].AcademicYear {
			return semester[i].SemesterNumber < semester[j].SemesterNumber
		}
		return semester[i].AcademicYear < semester[j].AcademicYear
	})

	//Calculate semester gpa and current gpa
	gpaResponse := response.GpaResponse{
		Semesters: semester,
	}
	totalCredits := 0
	for i := range semester {
		semester[i].Gpa /= float64(semester[i].TotalNumberOfCredits)
		gpaResponse.CurrentGpa += semester[i].Gpa * float64(semester[i].TotalNumberOfCredits)
		totalCredits += semester[i].TotalNumberOfCredits
	}
	gpaResponse.CurrentGpa = gpaResponse.CurrentGpa / float64(totalCredits)

	//Get future gpa
	if len(semester) >= 4 {
		gpaResponse.PredictedGpa = g.GetPredictedGpa(ctx,
			[]float64{semester[len(semester)-4].Gpa, semester[len(semester)-3].Gpa, semester[len(semester)-2].Gpa, semester[len(semester)-1].Gpa})
	} else {
		//sum := float64(0)
		//for i := range semester {
		//	sum += semester[i].Gpa
		//}
		//gpaResponse.PredictedGpa = sum / float64(len(semester))

		var semesterResult []float64
		for i := range semester {
			semesterResult = append(semesterResult, semester[i].Gpa)
		}
		for {
			if len(semesterResult) >= 4 {
				break
			}
			semesterResult = append(semesterResult, semester[len(semester)-1].Gpa)
		}
		gpaResponse.PredictedGpa = g.GetPredictedGpa(ctx, semesterResult)
	}
	return gpaResponse, nil
}

func (g *gpaService) GetPredictedGpa(ctx context.Context, gpa []float64) float64 {
	req := request.GetPredictedGpaRequest{
		GpaT1: gpa[0],
		GpaT2: gpa[1],
		GpaT3: gpa[2],
		GpaT4: gpa[3],
	}
	res, err := g.client.GetPredictedGpa(ctx, req)
	if err != nil {
		log.Println("GetPredictedGpa err: ", err)
		return (gpa[0] + gpa[1] + gpa[2] + gpa[3]) / 4
	}
	r := res.(response.GetPredictedGpaResponse)
	if len(r.Error) > 0 {
		log.Println("GetPredictedGpa err: ", r.Error)
		return (gpa[0] + gpa[1] + gpa[2] + gpa[3]) / 4
	}
	return r.PredictedFutureGpa
}

func NewGpaService(gpaRepo repo.GpaRepo, client *client.ClientEndpoint, redisGpaRepo repo.RedisGpaRepo) GpaService {
	return &gpaService{
		gpaRepo:      gpaRepo,
		client:       client,
		redisGpaRepo: redisGpaRepo,
	}
}
