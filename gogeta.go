package main;

import (
    // "github.com/herman-rogers/KingKai"
    "encoding/json"
    "log"
    "net/http"
    "golang.org/x/net/context"
    httptransport "github.com/go-kit/kit/transport/http"
);
//https://github.com/go-kit/kit.git
func main() {
    // kingkai.StartKingKai(routes);
    ctx := context.Background();
    svc := stringService{};

    countHandler := httptransport.NewServer(
        ctx,
        makeCountEndpoint(svc),
        decodeCountRequest,
        encodeResponse,
    );
    http.Handle("/count", countHandler);
    log.Fatal(http.ListenAndServe(":9000", nil));
}

func decodeCountRequest(r *http.Request) (interface{}, error) {
    var request countRequest;
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        return nil, err;
    }
    return request, nil;
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response);
}
