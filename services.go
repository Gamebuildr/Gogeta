package main;

type GogetaServiceInterface interface {
    Count(string) int;
}

type gogetaService struct{};
