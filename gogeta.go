package main;

import (
    "github.com/herman-rogers/kingkai"
    "github.com/hudl/fargo"
);

func main() {
    kingkai.StartKingKai(routes, ":9000");
}

func RegisterEureka() {
    RegisterEureka();
    e, _ := fargo.NewConnFromConfigFile("/etc/fargo.gcfg");
    e.AppWatchChannel
}
