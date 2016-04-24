package main;

import (
    "github.com/herman-rogers/kingkai"
    //"github.com/hudl/fargo"
);

func main() {
    RegisterEureka();
    kingkai.StartKingKai(routes, "");
}

func RegisterEureka() {
    //e := fargo.NewConn("http://eureka-gamebuildr.herokuapp.com");
    // // app, _ := e.GetApp("TESTAPP");
    //e.GetApps();
    // fmt.Println(apps);
    // for k, v := range apps {
    //     fmt.Println("k:", k, "v:", v);
    // }

    // e, _ := fargo.NewConnFromConfigFile("/etc/fargo.gcfg")
    // app, _ := e.GetApp("TESTAPP")
    // // starts a goroutine that updates the application on poll interval
    // e.UpdateApp(&app)
    // for {
    //     for _, ins := range app.Instances {
    //         fmt.Printf("%s, ", ins.HostName)
    //     }
    //     fmt.Println(len(app.Instances))
    //     <-time.After(10 * time.Second)
    // }
}
