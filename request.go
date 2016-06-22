package main

type gogetaService struct{}

type GogetaServiceInterface interface {
	GitClone(gitServiceRequest) string
	GitFindRepo(gitServiceRequest) string
}

type gitServiceRequest struct {
	Usr     string `json:"usr"`
	Repo    string `json:"repo"`
	Project string `json:"project"`
}

type serviceResponse struct {
	res string `json:"r"`
}
