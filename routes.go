package main;

import (
    "net/http"
    "encoding/json"
    "golang.org/x/net/context"
    httptransport "github.com/go-kit/kit/transport/http"
);

func routes() {
    context := context.Background();
    service := gogetaService{};

    http.Handle("/0/gitclone", gitCloneServerRequest(context, service));
}

func gitCloneServerRequest(context context.Context, service gogetaService) http.Handler {
    return httptransport.NewServer(
        context,
        makeGitCloneEndpoint(service),
        decodeGitCloneRequest,
        encodeResponse,
    );
}

func decodeGitCloneRequest(r *http.Request) (interface{}, error) {
    var request gitCloneRequest;
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        return nil, err;
    }
    return request, nil;
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response);
}
