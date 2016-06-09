package main

type GogetaServiceInterface interface {
	GitClone(gitServiceRequest) string
}

type gogetaService struct{}

type gitServiceRequest struct {
	Usr  string `json:"usr"`
	Repo string `json:"repo"`
}

type serviceResponse struct {
	res string `json:"r"`
}
