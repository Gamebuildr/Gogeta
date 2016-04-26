package main;

import (
    "log"
    "net/http"
    "os"
);

func main() {
    var port string = GetPort();
    routes();
    log.Fatal(http.ListenAndServe(port, nil));
}

func GetPort() string {
    var port = os.Getenv("PORT");
    if (port == "") {
        port = "9000";
        log.Printf("INFO: No PORT environment variable found, setting default.");
    }
    return ":" + port;
}
