package main

import (
    "../airlineproto"
    "../config"
    "../coordproto"
    "flag"
    "fmt"
    "strings"
    "strconv"
    // "log"
    // "net"
    "net/rpc"
    // "net/http"
    // "os"
    // "time"
)

type FlightStruct struct {
    FlightID string
    DepartureTime int64
    ArrivalTime int64
    DeparturePort string
    ArrivalPort string
    AvailableTickets int
    Capacity int
}

func parseFlight(str string) *airlineproto.FlightStruct {
    f := &airlineproto.FlightStruct{}
    parts := strings.Split(str, ":")
    f.FlightID = parts[0]
    f.DepartureTime, _ = strconv.ParseInt(parts[1], 10, 64)
    f.ArrivalTime, _ = strconv.ParseInt(parts[2], 10, 64)
    f.DeparturePort = parts[3]
    f.ArrivalPort = parts[4]
    tmp, _ := strconv.ParseInt(parts[5], 10, 32)
    f.Capacity = int(tmp)
    f.AvailableTickets = f.Capacity
    return f
}

var configFile *string = flag.String("config", "config/config", "configuration file path")
var numTimes *int = flag.Int("n", 1, "Number of times to execute the get or put.")

func main() {
    flag.Parse()

    config, err := config.ReadConfigFile(*configFile)

    serverRPC := make(map[string] *rpc.Client)
    for key, value := range(config.AirlineAddr) {
        serverRPC[key], err = rpc.DialHTTP("tcp", value)
        if err != nil {
            fmt.Printf("Could not connect to server %s, returning nil\n", value)
            return
        }
    }
    coordRPC, _ := rpc.DialHTTP("tcp", config.CoordAddr)

    cmd := flag.Arg(0)
    for et:=0; et < *numTimes; et++ {
        switch(cmd){
        case "a":
            flightStr := flag.Arg(1)
            flight := parseFlight(flightStr)
            args := &airlineproto.AddArgs{}
            var reply airlineproto.AddReply
            args.Flight = *flight
            airlineserver := strings.Split(flight.FlightID, "-")[0]
            serverRPC[airlineserver].Call("AirlineServerRPC.AddFlight", args, &reply)
            fmt.Println(*args, reply)
        case "d":
            args := &airlineproto.DeleteArgs{}
            var reply airlineproto.DeleteReply
            args.FlightID = flag.Arg(1)
            airlineserver := strings.Split(args.FlightID, "-")[0]
            serverRPC[airlineserver].Call("AirlineServerRPC.DeleteFlight", args, &reply)
            fmt.Println(*args, reply)
        case "q":
            args := &airlineproto.QueryArgs{}
            parts := strings.Split(flag.Arg(1), ":")
            args.StartTime, _ = strconv.ParseInt(parts[0], 10, 64)
            args.EndTime, _ = strconv.ParseInt(parts[1], 10, 64)
            for _, value := range(serverRPC) {
                var reply airlineproto.QueryReply
                value.Call("AirlineServerRPC.QueryFlights", args, &reply)
                fmt.Println("q", *args, reply)
            }
        case "r":
            flightStr := flag.Arg(1)
            flight := parseFlight(flightStr)
            args := &airlineproto.RescheduleArgs{}
            var reply airlineproto.RescheduleReply
            args.NewFlight = *flight
            args.OldFlightID = flight.FlightID
            airlineserver := strings.Split(flight.FlightID, "-")[0]
            serverRPC[airlineserver].Call("AirlineServerRPC.RescheduleFlight", args, &reply)
            fmt.Println(*args, reply)
        case "b":
            argStr := flag.Arg(1)
            parts := strings.Split(argStr, ":")
            args := &coordproto.BookArgs{}
            args.Email = parts[0]
            args.Count, _ = strconv.Atoi(parts[1])
            flightCnt, _ := strconv.Atoi(parts[2])
            for i:=0; i<flightCnt; i++ {
                flightID := parts[i+3]
                args.Flights = append(args.Flights, flightID)
            }
            fmt.Println(args.Flights)
            var reply coordproto.BookReply
            coordRPC.Call("CoordinatorRPC.BookFlights", args, &reply)
            fmt.Println("b", *args, reply)
        case "c":
            argStr := flag.Arg(1)
            parts := strings.Split(argStr, ":")
            args := &coordproto.BookArgs{}
            args.Email = parts[0]
            args.Count, _ = strconv.Atoi(parts[1])
            flightCnt, _ := strconv.Atoi(parts[2])
            for i:=0; i<flightCnt; i++ {
                flightID := parts[i+3]
                args.Flights = append(args.Flights, flightID)
            }
            fmt.Println(args.Flights)
            var reply coordproto.BookReply
            coordRPC.Call("CoordinatorRPC.CancelFlights", args, &reply)
            fmt.Println("c", *args, reply)
        }
    }

    for _, value := range(serverRPC) {
        value.Close()
    }
    coordRPC.Close()

    fmt.Println("Libclient Finished")
}
