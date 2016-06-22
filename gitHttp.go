package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

func gitCloneServiceRequest(context context.Context, service gogetaService) http.Handler {
	return httptransport.NewServer(
		context,
		makeGitCloneEndpoint(service),
		decodeGitRequest,
		gitEncodeResponse,
	)
}

func gitFindRepoServiceRequest(context context.Context, service gogetaService) http.Handler {
	return httptransport.NewServer(
		context,
		makeFindRepoEndpoint(service),
		decodeGitRequest,
		gitEncodeResponse,
	)
}

func makeGitCloneEndpoint(service GogetaServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(gitServiceRequest)
		res := service.GitClone(req)
		return serviceResponse{res}, nil
	}
}

func makeFindRepoEndpoint(service GogetaServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(gitServiceRequest)
		res := service.GitFindRepo(req)
		return serviceResponse{res}, nil
	}
}

func decodeGitRequest(r *http.Request) (interface{}, error) {
	var request gitServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func gitEncodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
