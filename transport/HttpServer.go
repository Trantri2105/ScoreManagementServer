package transport

import (
	"ScoreManagementSystem/client"
	"ScoreManagementSystem/dto/request"
	"ScoreManagementSystem/dto/response"
	endpoints "ScoreManagementSystem/endpoint"
	"ScoreManagementSystem/middleware"
	"ScoreManagementSystem/model"
	"ScoreManagementSystem/repo"
	"ScoreManagementSystem/service"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

var authMiddleware *middleware.Middleware

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(response.Message{Error: err.Error()})
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req model.Student
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeRegisterResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req request.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetStudentInfo(_ context.Context, r *http.Request) (interface{}, error) {
	studentId, err := authMiddleware.AuthenticateUserAndExtractUserId(r)
	if err != nil {
		return nil, err
	}
	return studentId, nil
}

func encodeGetStudentInfoResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetStudentGpaRequest(_ context.Context, r *http.Request) (interface{}, error) {
	studentId, err := authMiddleware.AuthenticateUserAndExtractUserId(r)
	if err != nil {
		return nil, err
	}
	return studentId, nil
}

func encodeGetStudentGpaRequestResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}
func NewHTTPServer(db *sql.DB, redisClient *redis.Client) *gin.Engine {
	c, err := client.NewClientEndpoint("localhost:5000")
	if err != nil {
		log.Fatal()
	}

	studentRepo := repo.NewStudentRepo(db)
	gpaRepo := repo.NewGpaRepo(db)

	redisStudentRepo := repo.NewRedisStudentRepo(redisClient)
	redisGpaRepo := repo.NewRedisGpaRepo(redisClient)

	jwtService := service.NewJwtService()
	studentService := service.NewStudentService(studentRepo, jwtService, redisStudentRepo)
	gpaService := service.NewGpaService(gpaRepo, c, redisGpaRepo)

	studentEndpoint := endpoints.NewStudentEndpoint(studentService)
	gpaEndpoint := endpoints.NewGpaEndpoint(gpaService)

	authMiddleware = middleware.NewMiddleware(jwtService)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	registerHandler := httptransport.NewServer(
		studentEndpoint.Register(),
		decodeRegisterRequest,
		encodeRegisterResponse,
		options...,
	)

	loginHandler := httptransport.NewServer(
		studentEndpoint.Login(),
		decodeLoginRequest,
		encodeLoginResponse,
		options...,
	)

	getStudentInfoHandler := httptransport.NewServer(
		studentEndpoint.GetStudentInfo(),
		decodeGetStudentInfo,
		encodeGetStudentInfoResponse,
		options...,
	)

	getStudentGpaEndpoint := httptransport.NewServer(
		gpaEndpoint.GetStudentGpa(),
		decodeGetStudentGpaRequest,
		encodeGetStudentInfoResponse,
		options...,
	)

	r := gin.Default()

	authRoute := r.Group("/auth")
	authRoute.POST("/register", gin.WrapH(registerHandler))
	authRoute.POST("/login", gin.WrapH(loginHandler))

	studentRoute := r.Group("/student")
	studentRoute.GET("/info", gin.WrapH(getStudentInfoHandler))

	gpaRoute := r.Group("/gpa")
	gpaRoute.GET("/detail", gin.WrapH(getStudentGpaEndpoint))

	return r
}
