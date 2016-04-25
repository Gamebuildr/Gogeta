package main;

// interface
type StringService interface {
    Count(string) int;
}

//implementation
type stringService struct{};

func (stringService) Count(s string) int {
    return len(s);
}
