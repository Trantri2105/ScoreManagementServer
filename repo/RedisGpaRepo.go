package repo

import (
	"ScoreManagementSystem/dto/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisGpaRepo interface {
	GetStudentGpaByStudentId(ctx context.Context, studentId string) (response.GpaResponse, error)
	SaveStudentGpa(ctx context.Context, studentId string, gpaResponse response.GpaResponse)
}

type redisGpaRepo struct {
	redisClient *redis.Client
}

func (r *redisGpaRepo) GetStudentGpaByStudentId(ctx context.Context, studentId string) (response.GpaResponse, error) {
	gpaKey := fmt.Sprintf("gpa#%s", studentId)
	result, err := r.redisClient.Get(ctx, gpaKey).Result()
	if err != nil {
		log.Println("Redis gpa repo, get student gpa from redis err: ", err)
		return response.GpaResponse{}, err
	}
	var gpaRes response.GpaResponse
	err = json.Unmarshal([]byte(result), &gpaRes)
	if err != nil {
		log.Println("Redis gpa repo, unmarshal student gpa cache err: ", err)
		return response.GpaResponse{}, err
	}
	return gpaRes, nil
}

func (r *redisGpaRepo) SaveStudentGpa(ctx context.Context, studentId string, gpaResponse response.GpaResponse) {
	gpaKey := fmt.Sprintf("gpa#%s", studentId)
	gpaBytes, err := json.Marshal(gpaResponse)
	if err != nil {
		log.Println("Redis gpa repo, marshal gpa response err: ", err)
	} else {
		_, err = r.redisClient.Set(ctx, gpaKey, gpaBytes, 0).Result()
		if err != nil {
			log.Println("Redis gpa repo, save gpa to redis err: ", err)
		}
	}
}

func NewRedisGpaRepo(redisClient *redis.Client) RedisGpaRepo {
	return &redisGpaRepo{redisClient: redisClient}
}
