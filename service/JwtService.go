package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"time"
)

type JwtService interface {
	CreateToken(studentId string) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type jwtService struct{}

func (*jwtService) CreateToken(studentId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"studentId": studentId,
		"exp":       time.Now().Add(60 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println("Jwt service, create access token err :", err)
		return "", errors.New("internal server error")
	}
	return tokenString, nil
}

func (*jwtService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

func NewJwtService() JwtService {
	return &jwtService{}
}
