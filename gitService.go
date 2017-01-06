package main

import (
    "errors"
    "fmt"
    "os/exec"

    "github.com/Gamebuildr/Gogeta/config"
    "github.com/Gamebuildr/Gogeta/logger"
    uuid "github.com/satori/go.uuid"
    git "gopkg.in/libgit2/git2go.v23"
)

const UP_TO_DATE = "UP_TO_DATE"
const NORMAL_MERGE = "NORMAL_MERGE"
const FAST_FORWARD = "FAST_FORWARD"

func GitProcessSQSMessages(data scmServiceRequest) error {
    if data.Usr == "" || data.Project == "" || data.Repo == "" {
        message := "Missing git data"
        SendRawMessage(
            data.Id,
            data.Buildcount,
            message,
            message)
        return errors.New("Missing Git Service Properties")
    }
    switch data.Type {
    case "GOGETA_CLONE":
        go GitShallowClone(data)
    case "GOGETA_NEW_BUILD":
        go RunNewGitBuild(data)
    }
    return nil
}

func GitShallowClone(data scmServiceRequest) {
    repoPath, relativePath := GetRepoPath(data.Project)
    cmd := exec.Command("git", "clone", "--depth", "2", data.Repo, repoPath)

    cloneMsgDev := "git clone --depth 2 " + data.Repo + " " + repoPath
    SendRawMessage(
        data.Id,
        data.Buildcount,
        "git clone started",
        cloneMsgDev)

    logfile := logger.GetLogFile()
    defer logfile.Close()

    cmd.Stdout = logfile
    cmd.Stderr = logfile

    commandErr := cmd.Start()
    logger.LogError(commandErr, "Git Clone")
    if commandErr != nil {
        SendRawMessage(
            data.Id,
            data.Buildcount,
            "git clone failed",
            commandErr.Error())
    }

    cloneErr := cmd.Wait()
    logger.LogError(cloneErr, "Git Clone")
    if cloneErr != nil {
        SendRawMessage(
            data.Id,
            data.Buildcount,
            "git clone failed",
            cloneErr.Error())
    }
    if cloneErr == nil {
        cloneSuccess := "git clone succeeded"
        SendRawMessage(
            data.Id,
            data.Buildcount,
            cloneSuccess,
            cloneSuccess)
        gitData := &GogetaRepo{
            data.Id,
            data.Usr,
            data.Repo,
            relativePath,
            data.SCMType,
            data.Engine,
            data.Platform,
            data.Buildcount,
        }
        CreateGitCredentials(repoPath)
        go SaveRepo(*gitData)
        go TriggerMrRobotBuild(*gitData)
    }
}

func RunNewGitBuild(data scmServiceRequest) {
    gitData := FindRepo(data.Usr, data.Id)
    gitData.BuildCount = gitData.BuildCount + 1
    // put raw message in new build
    SendRawMessage(
        gitData.BuildrId,
        gitData.BuildCount,
        "new build started",
        "new build started")
    go TriggerMrRobotBuild(gitData)
    go UpdateRepo(gitData)
}

func UpdateGitRepositories() {
    repos := FindAllRepos()
    for i := 0; i < len(repos); i++ {
        repoPath, err := config.MainConfig.GetConfigKey("RepoPath")
        if err != nil {
            fmt.Printf(err.Error())
        }
        folder := repoPath + repos[i].Folder
        repo, err := git.OpenRepository(folder)
        if err != nil {
            logger.Error(err.Error())
            continue
        }
        if repo.IsBare() {
            return
        }
        go GitPull(repos[i], repo)
    }
}

func GetRepoPath(project string) (string, string) {
    var uuid = uuid.NewV4()
    relativePath := project + "_" + uuid.String()
    repoPath, err := config.MainConfig.GetConfigKey("RepoPath")
    if err != nil {
        fmt.Printf(err.Error())
    }
    fullPath := repoPath + relativePath
    return fullPath, relativePath
}

func CreateGitCredentials(repo string) {
    openRepo, err := git.OpenRepository(repo)
    if err != nil {
        logger.LogError(err, "Git Create Credentials")
    }
    config, configErr := openRepo.Config()
    if configErr != nil {
        logger.LogError(configErr, "Git Config Set")
    }
    configNameErr := config.SetString("user.name", "gamebuildr")
    if configNameErr != nil {
        logger.LogError(configNameErr, "Git Config Name")
    }
    configEmailErr := config.SetString("user.email", "contact@gamebuildr.io")
    if configEmailErr != nil {
        logger.LogError(configEmailErr, "Git Config Email")
    }
}

func GitPull(gitData GogetaRepo, repo *git.Repository) {
    FetchLatestsUpdates(gitData)
    remoteBranch := GetRemoteBranch(repo)
    remoteBranchID := remoteBranch.Target()

    analysis, err := AnalyzeLatestUpdates(repo, remoteBranch)
    if err != nil {
        logger.Error(err.Error())
    }
    mergeType := CheckMergeAnalysis(analysis)
    msg := "Git Merge " + mergeType

    switch mergeType {
    case UP_TO_DATE:
        break
    case NORMAL_MERGE:
        err := MergeNormal(repo, remoteBranch, remoteBranchID)
        gitData.BuildCount = gitData.BuildCount + 1
        //put raw message in new build
        SendRawMessage(
            gitData.BuildrId,
            gitData.BuildCount,
            "new code changes detected",
            "git merge origin master")
        BuildAfterMerge(err, msg, gitData)
        go UpdateRepo(gitData)
        break
    case FAST_FORWARD:
        err := MergeFastForward(repo, remoteBranchID)
        gitData.BuildCount = gitData.BuildCount + 1
        //put raw message in new build
        SendRawMessage(
            gitData.BuildrId,
            gitData.BuildCount,
            "new code changes detected",
            "git merge fast forward")
        BuildAfterMerge(err, msg, gitData)
        go UpdateRepo(gitData)
        break
    }
}

func SendRawMessage(dataId string, dataCount int, message string, devMessage string) {
    gitMessage := GamebuildrMessage{
        dataId,
        dataCount,
        message,
        devMessage,
        "BUILDR_MESSAGE",
    }
    SendGamebuildrMessage(gitMessage)
}

func GetRemoteBranch(repo *git.Repository) *git.Reference {
    remoteBranch, err := repo.References.Lookup("refs/remotes/origin/master")
    logger.LogError(err, "Git Remote Branch")
    return remoteBranch
}

func FetchLatestsUpdates(data GogetaRepo) {
    cmd := exec.Command("git", "fetch", "--depth", "1", "origin", "master")
    repoPath, err := config.MainConfig.GetConfigKey("RepoPath")
    if err != nil {
        fmt.Printf(err.Error())
    }
    cmd.Dir = repoPath + data.Folder

    commandErr := cmd.Start()
    logger.LogError(commandErr, "Git Fetch")

    fetchErr := cmd.Wait()
    logger.LogError(fetchErr, "Git Fetch")
}

func AnalyzeLatestUpdates(repo *git.Repository, remoteBranch *git.Reference) (git.MergeAnalysis, error) {
    annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
    if err != nil {
        var nullAnalysis git.MergeAnalysis
        return nullAnalysis, err
    }
    mergeHeads := make([]*git.AnnotatedCommit, 1)
    mergeHeads[0] = annotatedCommit
    analysis, _, err := repo.MergeAnalysis(mergeHeads)
    logger.LogError(err, "Git Merge Analysis")
    return analysis, nil
}

func CheckMergeAnalysis(analysis git.MergeAnalysis) string {
    if analysis&git.MergeAnalysisUpToDate != 0 {
        return UP_TO_DATE
    }
    if analysis&git.MergeAnalysisNormal != 0 {
        return NORMAL_MERGE
    }
    return FAST_FORWARD
}

func MergeNormal(repo *git.Repository, remoteBranch *git.Reference, remoteBranchID *git.Oid) error {
    head, err := repo.Head()
    if err != nil {
        return err
    }
    annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
    if err != nil {
        return err
    }
    if err := repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
        return err
    }
    index, err := repo.Index()
    if err != nil {
        return err
    }
    if index.HasConflicts() {
        return errors.New("Conflicts encountered, Please Resolve Them")
    }
    sig, err := repo.DefaultSignature()
    if err != nil {
        return err
    }
    treeId, err := index.WriteTree()
    if err != nil {
        return err
    }
    tree, err := repo.LookupTree(treeId)
    if err != nil {
        return err
    }
    localCommit, err := repo.LookupCommit(head.Target())
    if err != nil {
        return err
    }
    remoteCommit, err := repo.LookupCommit(remoteBranchID)
    if err != nil {
        return err
    }
    repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)
    repo.StateCleanup()
    return nil
}

func MergeFastForward(repo *git.Repository, remoteBranchID *git.Oid) error {
    head, err := repo.Head()
    if err != nil {
        return err
    }
    remoteTree, err := repo.LookupTree(remoteBranchID)
    if err != nil {
        return err
    }
    if err := repo.CheckoutTree(remoteTree, nil); err != nil {
        return err
    }
    branchRef, err := repo.References.Lookup("refs/heads/master")
    if err != nil {
        return err
    }
    branchRef.SetTarget(remoteBranchID, "")
    if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
        return err
    }
    return nil
}
