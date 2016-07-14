package main

import (
	"errors"
	"github.com/herman-rogers/gogeta/config"
	"github.com/herman-rogers/gogeta/logger"
	git "github.com/libgit2/git2go"
	"github.com/satori/go.uuid"
	"os/exec"
)

const UP_TO_DATE = "UP_TO_DATE"
const NORMAL_MERGE = "NORMAL_MERGE"
const FAST_FORWARD = "FAST_FORWARD"

func GitProcessSQSMessages(gitReq gitServiceRequest) error {
	if gitReq.Usr == "" || gitReq.Project == "" || gitReq.Repo == "" {
		return errors.New("Missing Git Service Properties")
	}
	go GitShallowClone(gitReq)
	return nil
}

func GitShallowClone(data gitServiceRequest) {
	var uuid uuid.UUID = uuid.NewV4()
	folder := config.File.RepoPath + data.Usr + "/" + data.Project + "_" + uuid.String()
	cmd := exec.Command("git", "clone", "--depth", "1", data.Repo, folder)

	logfile := logger.GetLogFile()
	defer logfile.Close()

	cmd.Stdout = logfile
	cmd.Stderr = logfile

	commandErr := cmd.Start()
	logger.LogError(commandErr, "Git Clone")

	cloneErr := cmd.Wait()
	logger.LogError(cloneErr, "Git Clone")
	if cloneErr == nil {
		gitData := &GogetaRepo{data.Usr, data.Repo, folder}
		go SaveRepo(gitData)
	}
}

func UpdateGitRepositories() {
	repos := FindAllRepos()
	for i := 0; i < len(repos); i++ {
		folder := repos[i].Folder
		repo, err := git.OpenRepository(folder)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		if repo.IsBare() {
			return
		}
		go GitPull(repos[i], repo)
	}
}

func GitPull(data GogetaRepo, repo *git.Repository) {
	FetchLatestsUpdates(data)
	analysis := AnalyzeLatestUpdates(repo)
	mergeType := CheckMergeAnalysis(analysis)
	var mergeErr error
	switch mergeType {
	case UP_TO_DATE:
		break
	case NORMAL_MERGE:
		mergeErr = MergeNormal(repo)
	case FAST_FORWARD:
		mergeErr = MergeFastForward(repo)
	}
	logger.LogError(mergeErr, "Merge Error")
}

func FetchLatestsUpdates(data GogetaRepo) {
	cmd := exec.Command("git", "fetch", "--depth", "1", "origin", "master")
	cmd.Dir = data.Folder

	commandErr := cmd.Start()
	logger.LogError(commandErr, "Git Fetch")

	cloneErr := cmd.Wait()
	logger.LogError(cloneErr, "Git Fetch")
}

func AnalyzeLatestUpdates(repo *git.Repository) git.MergeAnalysis {
	remoteBranch := GetRemoteBranch(repo)
	annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		logger.LogError(err, "Git Merge Analysis")
	}
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = annotatedCommit
	analysis, _, err := repo.MergeAnalysis(mergeHeads)
	logger.LogError(err, "Git Merge Analysis")
	return analysis
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

func MergeNormal(repo *git.Repository) error {
	head, err := repo.Head()
	if err != nil {
		return err
	}
	remoteBranch := GetRemoteBranch(repo)
	remoteBranchID := remoteBranch.Target()
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

func MergeFastForward(repo *git.Repository) error {
	head, err := repo.Head()
	if err != nil {
		return err
	}
	remoteBranch := GetRemoteBranch(repo)
	remoteBranchID := remoteBranch.Target()
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

func GetRemoteBranch(repo *git.Repository) *git.Reference {
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/master")
	logger.LogError(err, "Git Remote Branch")
	return remoteBranch
}
