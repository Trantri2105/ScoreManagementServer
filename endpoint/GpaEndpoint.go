package endpoint

import (
	"ScoreManagementSystem/service"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type GpaEndpoint interface {
	GetStudentGpa() endpoint.Endpoint
}

type gpaEndpoint struct {
	gpaService service.GpaService
}

func (g *gpaEndpoint) GetStudentGpa() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		return g.gpaService.GetStudentGpaByStudentId(ctx, req)
	}
}

func NewGpaEndpoint(gpaService service.GpaService) GpaEndpoint {
	return &gpaEndpoint{gpaService: gpaService}
}
