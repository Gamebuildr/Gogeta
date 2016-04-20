package main;

import "net/http";

import (
    "fmt"
    "html"
);

type Route struct {
    Name string
    Method string
    Pattern string
    HandlerFunc http.HandlerFunc
}

type Routes []Route;

var routes = Routes {
    Route {
        "Example",
        "GET",
        "/",
        Index,
    },
    Route {
        "RouterExample",
        "GET",
        "/example",
        ExampleRoute,
    },
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path));
}

func ExampleRoute(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Example Route, %q", html.EscapeString(r.URL.Path));
}
