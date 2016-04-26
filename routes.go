package main;

import (
    "net/http"
    "encoding/json"
);

func routes() {
    http.Handle("/count", countServerRequest());
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
