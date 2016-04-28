package main;

import (
    "log"
    "golang.org/x/net/context"
    "github.com/go-kit/kit/endpoint"
    "github.com/satori/go.uuid"
    git "github.com/libgit2/git2go"
);

func makeGitCloneEndpoint(service GogetaServiceInterface) endpoint.Endpoint {
    return func (ctx context.Context, request interface{}) (interface{}, error)  {
        req := request.(gitServiceRequest);
        res := service.GitClone(req);
        return serviceResponse{res}, nil;
    }
}

func (gogetaService) GitClone(gitReq gitServiceRequest) string {
    go StartGitClone(gitReq.Repo)
    return "Clone Success"
}

func StartGitClone(gitRepo string) {
    var cloneOptions *git.CloneOptions = &git.CloneOptions{}
    folder := uuid.NewV4()
    repo, err := git.Clone(gitRepo, folder.String(), cloneOptions)
    if(err != nil) {
        log.Print("ERROR: " + err.Error());
    }
    log.Print(repo)
}
