package main;

import (
    "golang.org/x/net/context"
    "github.com/go-kit/kit/endpoint"
    httptransport "github.com/go-kit/kit/transport/http"
    "net/http"
);

// Initial
func countServerRequest() http.Handler {
    context := context.Background();
    service := gogetaService{};

    return httptransport.NewServer(
        context,
        makeCountEndpoint(service),
        decodeCountRequest,
        encodeResponse,
    );
}

// Service Endpoint
func makeCountEndpoint(service GogetaServiceInterface) endpoint.Endpoint {
    return func (ctx context.Context, request interface{}) (interface{}, error)  {
        req := request.(countRequest);
        v := service.Count(req.S);
        return countResponse{v}, nil;
    }
}

func (gogetaService) Count(s string) int {
    return len(s);
}
