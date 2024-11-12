package middleware

import (
	"ScoreManagementSystem/service"
	"errors"
	"net/http"
	"strings"
)

type Middleware struct {
	jwtService service.JwtService
}

func (m Middleware) AuthenticateUserAndExtractUserId(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", errors.New("require access token")
	}
	header := strings.Fields(authHeader)
	if len(header) != 2 && header[0] != "Bearer" {
		return "", errors.New("wrong access token format")
	}

	accessToken := header[1]
	claims, err := m.jwtService.VerifyToken(accessToken)
	if err != nil {
		return "", err
	}
	return claims["studentId"].(string), nil
}

func NewMiddleware(jwtService service.JwtService) *Middleware {
	return &Middleware{jwtService: jwtService}
}
