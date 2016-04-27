package main;

import (
    "golang.org/x/net/context"
    "github.com/go-kit/kit/endpoint"
    //"github.com/libgit2/git2go"
);

func makeGitCloneEndpoint(service GogetaServiceInterface) endpoint.Endpoint {
    return func (ctx context.Context, request interface{}) (interface{}, error)  {
        req := request.(gitCloneRequest);
        r := service.GitClone(req.clone);
        return serviceResponse{r}, nil;
    }
}

func (gogetaService) GitClone(s string) string {
    return "Git Clone";
}
