package main;

import (
    "net/http"
    "encoding/json"
    "golang.org/x/net/context"
    "github.com/go-kit/kit/endpoint"
    httptransport "github.com/go-kit/kit/transport/http"
);

func gitCloneServerRequest(context context.Context, service gogetaService) http.Handler {
    return httptransport.NewServer(
        context,
        makeGitCloneEndpoint(service),
        decodeGitRequest,
        gitEncodeResponse,
    );
}

func makeGitCloneEndpoint(service GogetaServiceInterface) endpoint.Endpoint {
    return func (ctx context.Context, request interface{}) (interface{}, error)  {
        req := request.(gitServiceRequest);
        res := service.GitClone(req);
        return serviceResponse{res}, nil;
    }
}

func decodeGitRequest(r *http.Request) (interface{}, error) {
    var request gitServiceRequest;
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        return nil, err;
    }
    return request, nil;
}

func gitEncodeResponse(w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response);
}
