package main;

//Request and Responses
type countRequest struct {
    S string `json:"s"`
}

type countResponse struct {
    V int `json:"v"`
}
