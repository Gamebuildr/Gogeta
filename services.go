package main;

type GogetaServiceInterface interface {
    GitClone(string) string;
}

type gogetaService struct{};

type gitCloneRequest struct {
    clone string `json:"clone"`
}

type serviceResponse struct {
    R string `json:"r"`
}
