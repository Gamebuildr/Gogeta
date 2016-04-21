package kingkai;

import (
    "log"
    "net/http"
);

func StartKingKai(routes Routes, port string) {
    router := CreateRouter(routes);
    log.Fatal(http.ListenAndServe(port, router));
}
