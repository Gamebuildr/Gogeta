package main;

import (
    //"github.com/herman-rogers/KingKai"
    "github.com/herman-rogers/Mr.Robot"
);

func main() {
    RegisterMrRobot();
    //kingkai.StartKingKai(routes);
}

func RegisterMrRobot() {
    mrrobot.SetLoggerAppName("gogeta");
    client := mrrobot.NewClient([]string {
        //"http://eureka-gamebuildr.herokuapp.com",
        "http://localhost:8080/eureka",
    });
    //instance := eureka.NewInstanceInfo("gogeta.herokuapp.com", "gogeta", "gogeta.herokuapp.com", 80, 30, false);
    instance := mrrobot.NewInstanceInfo("test.com", "test", "69.172.200.235", 80, 30, false);
    instance.Metadata = &mrrobot.MetaData {
        Map: make(map[string]string),
    }
    instance.Metadata.Map["foo"] = "bar";
    client.RegisterInstance("gogeta", instance);
    client.GetApplication(instance.App);
    client.GetInstance(instance.App, instance.HostName);
    client.SendHeartbeatUpdates(instance.App, instance.HostName);
}
