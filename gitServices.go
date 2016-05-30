package main;

import (
    "log"
    // "time"
    //"github.com/kardianos/osext"
    "golang.org/x/net/context"
    "github.com/go-kit/kit/endpoint"
    "github.com/satori/go.uuid"
    // git "github.com/libgit2/git2go"
    "os/exec"
    "os"
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
    cmd := exec.Command("git", "clone", "--depth", "1", repo, location)
    logfile, err := os.Create("gogeta.log")
    if(err != nil) {
        log.Print("Clone Error: " + err.Error())
    }
    defer logfile.Close()

    cmd.Stdout = logfile
    cmd.Stderr = logfile

    if err := cmd.Start(); err != nil {
        log.Print("Clone Error: " + err.Error())
    }

    // ticker := time.NewTicker(time.Second)
    // go func(ticker *time.Ticker) {
    //         //now := time.Now()
    //     for _ = range ticker.C {
    //         //log.Print(out)
    //         log.Print("Ticker Running")
    //         // go io.Copy(writer, stdout)
    //     }
    // }(ticker)
    log.Print("Starting Clone")
    cmd.Wait()
    log.Print("Finished Clone")
    //folderPath, err := osext.ExecutableFolder()
}

// func GitUpdateRepo(repo string) {
//
// }

// func GitDeleteRepo(repo string) {
//
// }
