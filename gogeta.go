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
    c := fargo.NewConn("http://eureka-gamebuildr.herokuapp.com")
    c.GetApps()
}
