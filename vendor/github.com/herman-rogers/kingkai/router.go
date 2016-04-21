package kingkai;

import (
    "net/http"
    "github.com/gorilla/mux"
);

type Route struct {
    Name string
    Method string
    Pattern string
    HandlerFunc http.HandlerFunc
}

type Routes []Route;

func CreateRouter(customRoutes Routes) *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range customRoutes {
        var handler http.Handler;
        handler = route.HandlerFunc;
        handler = Logger(handler, route.Name);
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler);
    }
    return router;
}
