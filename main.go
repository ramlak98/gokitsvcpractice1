package main

import (
	"errors"
	"strings"
	"context"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"log"
	"encoding/json"
)

type StringService interface {
	Uppercase(ctx context.Context, string2 string) (string, error)
	Count(ctx context.Context, string2 string) int
}

// Implement interface
type stringService struct{}

func (stringService) Uppercase (_ context.Context, str string) (string, error) {
	if str == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(str), nil
}

func (stringService) Count(_ context.Context, str string) int {
	return len(str)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// Requests and responses
type uppercaseRequest struct{
	S string `json:"s"`
}

type uppercaseResponse struct {
	V string `json:"v"`
	Err string `json:"err,omitempty"`
}

type countRequest struct{
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

// endpoints

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{})(interface{}, error){
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(ctx, req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (interface{}, error){
		req := request.(countRequest)
		v  := svc.Count(ctx, req.S)
		return countResponse{v}, nil
	}
}

func main(){
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
		)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error){
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}