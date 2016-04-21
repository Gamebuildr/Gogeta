package main;

import (
    "github.com/herman-rogers/kingkai"
    "github.com/hudl/fargo"
);

func main() {
    RegisterEureka();
    kingkai.StartKingKai(routes, ":9000");
}

func RegisterEureka() {
    // e, _ := fargo.NewConnFromConfigFile("/etc/fargo.gcfg");
    // e.AppWatchChannel
    c := fargo.NewConn("http://127.0.0.1:8080/eureka/v2")
    c.GetApps()
}
