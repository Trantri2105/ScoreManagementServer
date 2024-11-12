package repo

import (
	"ScoreManagementSystem/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisStudentRepo interface {
	GetStudentInfoByStudentId(ctx context.Context, studentId string) (model.Student, error)
	SaveStudentInfo(ctx context.Context, studentInfo model.Student)
}

type redisStudentRepo struct {
	redisClient *redis.Client
}

func (r *redisStudentRepo) GetStudentInfoByStudentId(ctx context.Context, studentId string) (model.Student, error) {
	studentKey := fmt.Sprintf("student#%s", studentId)
	result, err := r.redisClient.Get(ctx, studentKey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Println("Redis student repo, get student from redis err: ", err)
		}
		return model.Student{}, err
	}
	var student model.Student
	err = json.Unmarshal([]byte(result), &student)
	if err != nil {
		log.Println("Redis student repo, unmarshal student cache err: ", err)
		return model.Student{}, err
	}
	return student, nil
}

func (r *redisStudentRepo) SaveStudentInfo(ctx context.Context, studentInfo model.Student) {
	studentKey := fmt.Sprintf("student#%s", studentInfo.Id)
	studentBytes, err := json.Marshal(studentInfo)
	if err != nil {
		log.Println("Redis student repo, student marshall err: ", err)
	} else {
		_, err = r.redisClient.Set(ctx, studentKey, studentBytes, 10*time.Minute).Result()
		if err != nil {
			log.Println("Redis student repo, set student cache err: ", err)
		}
	}
}

func NewRedisStudentRepo(redisClient *redis.Client) RedisStudentRepo {
	return &redisStudentRepo{redisClient: redisClient}
}
