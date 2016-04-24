package main;

import (
    "fmt"
    "os"
    // "time"
    "github.com/herman-rogers/kingkai"
    "github.com/ArthurHlt/go-eureka-client/eureka"
);

func main() {
    var port string = GetPort();
    RegisterEureka();
    kingkai.StartKingKai(routes, port);
}

func GetPort() string {
    var port = os.Getenv("PORT");
    if (port == "") {
        port = "9000";
        fmt.Println("INFO: No PORT environment variable found, setting default.");
    }
    return ":" + port;
}

func RegisterEureka() {
    client := eureka.NewClient([]string {
        "http://eureka-gamebuildr.herokuapp.com",
    });
    instance := eureka.NewInstanceInfo("gogeta.herokuapp.com", "gogeta", "gogeta.herokuapp.com", 80, 30, false);
    instance.Metadata = &eureka.MetaData {
        Map: make(map[string]string),
    }
    instance.Metadata.Map["foo"] = "bar";
    client.RegisterInstance("gogeta", instance);
    client.GetApplication(instance.App);
    client.GetInstance(instance.App, instance.HostName);
    client.SendHeartbeat(instance.App, instance.HostName);
}

// func main() {
//
//     client := eureka.NewClient([]string{
//         "http://127.0.0.1:8761/eureka", //From a spring boot based eureka server
//         // add others servers here
//     })
//     instance := eureka.NewInstanceInfo("test.com", "test", "69.172.200.235", 80, 30, false) //Create a new instance to register
//     instance.Metadata = &eureka.MetaData{
//         Map: make(map[string]string),
//     }
//     instance.Metadata.Map["foo"] = "bar" //add metadata for example
//     client.RegisterInstance("myapp", instance) // Register new instance in your eureka(s)
//     applications, _ := client.GetApplications() // Retrieves all applications from eureka server(s)
//     client.GetApplication(instance.App) // retrieve the application "test"
//     client.GetInstance(instance.App, instance.HostName) // retrieve the instance from "test.com" inside "test"" app
//     client.SendHeartbeat(instance.App, instance.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
// }
