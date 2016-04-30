package main;

import (
    "log"
    // "sync"
    "golang.org/x/net/context"
    "github.com/go-kit/kit/endpoint"
    "github.com/satori/go.uuid"
    // git "github.com/libgit2/git2go"
    "os/exec"
);

func makeGitCloneEndpoint(service GogetaServiceInterface) endpoint.Endpoint {
    return func (ctx context.Context, request interface{}) (interface{}, error)  {
        req := request.(gitServiceRequest);
        res := service.GitClone(req);
        return serviceResponse{res}, nil;
    }
}

func (gogetaService) GitClone(gitReq gitServiceRequest) string {
    go GitShallowClone(gitReq.Repo)
    return "Clone Success"
}

func GitShallowClone(repo string) {
    var folder uuid.UUID = uuid.NewV4()
    var location string = "./" + folder.String()

    log.Print("Git Clone Request");
    cmd := exec.Command("git", "clone", "--depth", "1", repo, location)
    err := cmd.Run()
    if(err != nil) {
        log.Print("Git Clone Error: " + err.Error())
    } else {
        log.Print("Git Clone Success")
    }
}
