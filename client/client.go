package client

import (
	"ScoreManagementSystem/dto/response"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ClientEndpoint struct {
	GetPredictedGpa endpoint.Endpoint
}

func encodeGetPredictedGpaRequest(_ context.Context, r *http.Request, req interface{}) error {
	r.URL.Path = "/predict"
	r.Header.Set("Content-Type", "application/json")
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	return nil

}

func decodeGetPredictedGpaResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var req response.GetPredictedGpaResponse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func NewClientEndpoint(instance string) (*ClientEndpoint, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	u.Path = ""
	if err != nil {
		return &ClientEndpoint{}, err
	}
	return &ClientEndpoint{
		GetPredictedGpa: httptransport.NewClient("POST", u, encodeGetPredictedGpaRequest, decodeGetPredictedGpaResponse).Endpoint(),
	}, nil
}
