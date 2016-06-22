package main

import (
	"net/http"

	"golang.org/x/net/context"
)

func routes() {
	context := context.Background()
	service := gogetaService{}

	http.Handle("/0/gitclone", gitCloneServiceRequest(context, service))
	http.Handle("/0/gitfindrepo", gitFindRepoServiceRequest(context, service))
}
