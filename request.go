package main

type gitServiceRequest struct {
	Usr     string `json:"usr"`
	Repo    string `json:"repo"`
	Project string `json:"project"`
}

type serviceResponse struct {
	res string `json:"r"`
}
