package main;

import (
    "golang.org/x/net/context"
    "github.com/go-kit/kit/endpoint"
);

func makeCountEndpoint(service StringService) endpoint.Endpoint {
    return func (ctx context.Context, request interface{}) (interface{}, error)  {
        req := request.(countRequest);
        v := service.Count(req.S);
        return countResponse{v}, nil;
    }
}
