package main

import (
	"errors"
	"os/exec"

	"github.com/herman-rogers/gogeta/config"
	"github.com/herman-rogers/gogeta/logger"
	git "github.com/libgit2/git2go"
	"github.com/satori/go.uuid"
)

const UP_TO_DATE = "UP_TO_DATE"
const NORMAL_MERGE = "NORMAL_MERGE"
const FAST_FORWARD = "FAST_FORWARD"

func GitProcessSQSMessages(gitReq scmServiceRequest) error {
	if gitReq.Usr == "" || gitReq.Project == "" || gitReq.Repo == "" {
		return errors.New("Missing Git Service Properties")
	}
	go GitShallowClone(gitReq)
	return nil
}

func GitShallowClone(data scmServiceRequest) {
	repoPath, relativePath := GetRepoPath(data.Project)
	cmd := exec.Command("git", "clone", "--depth", "2", data.Repo, repoPath)

	cloneMsgDev := "git clone --depth 2 " + data.Repo + " " + repoPath
	SendGitMessage(data, "git clone started", cloneMsgDev)

	logfile := logger.GetLogFile()
	defer logfile.Close()

	cmd.Stdout = logfile
	cmd.Stderr = logfile

	commandErr := cmd.Start()
	logger.LogError(commandErr, "Git Clone")
	if commandErr != nil {
		SendGitMessage(data, "git clone failed", commandErr.Error())
	}

	cloneErr := cmd.Wait()
	logger.LogError(cloneErr, "Git Clone")
	if cloneErr != nil {
		SendGitMessage(data, "git clone failed", cloneErr.Error())
	}
	if cloneErr == nil {
		cloneSuccess := "git clone succeded"
		SendGitMessage(data, cloneSuccess, cloneSuccess)
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

func SendGitMessage(data scmServiceRequest, message string, devMessage string) {
	gitMessage := GamebuildrMessage{
		data.Id,
		data.Buildcount,
		message,
		devMessage,
		"BUILDR_MESSAGE",
	}
	SendGamebuildrMessage(gitMessage)
}

func GetRepoPath(project string) (string, string) {
	var uuid uuid.UUID = uuid.NewV4()
	relativePath := project + "_" + uuid.String()
	repoPath := config.File.RepoPath + relativePath
	return repoPath, relativePath
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

func UpdateGitRepositories() {
	repos := FindAllRepos()
	for i := 0; i < len(repos); i++ {
		folder := config.File.RepoPath + repos[i].Folder
		repo, err := git.OpenRepository(folder)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if repo.IsBare() {
			continue
		}
		go GitPull(repos[i], repo)
	}
}

func GitPull(data GogetaRepo, repo *git.Repository) {
	FetchLatestsUpdates(data)
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
		BuildAfterMerge(err, msg, data)
		break
	case FAST_FORWARD:
		err := MergeFastForward(repo, remoteBranchID)
		BuildAfterMerge(err, msg, data)
		break
	}
}

func GetRemoteBranch(repo *git.Repository) *git.Reference {
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/master")
	logger.LogError(err, "Git Remote Branch")
	return remoteBranch
}

func FetchLatestsUpdates(data GogetaRepo) {
	cmd := exec.Command("git", "fetch", "--depth", "1", "origin", "master")
	cmd.Dir = config.File.RepoPath + data.Folder

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
