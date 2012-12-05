package main

import(
    "flag"
    "../consensus"
    "time"
    "runtime"
)

var configFile *string = flag.String("config", "config/config", "configuration file path")
var id *int = flag.Int("id", 0, "Number of times to execute the get or put.")
var name *string = flag.String("name", "1", "configuration file path")

func main() {
    runtime.GOMAXPROCS(20)
    flag.Parse()
    consensus.NewPaxosEngine(*configFile, *name, *id)
    for {
        time.Sleep(time.Second * 100)
    }
}
