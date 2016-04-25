package main;

import "net/http";

import (
    "fmt"
    "html"
    "github.com/herman-rogers/KingKai"
);

var routes = kingkai.Routes {
    kingkai.Route {
        "Example",
        "GET",
        "/",
        Index,
    },
    kingkai.Route {
        "RouterExample",
        "GET",
        "/clone",
        ExampleRoute,
    },
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path));
}

func ExampleRoute(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Example Route, %q", html.EscapeString(r.URL.Path));
}
