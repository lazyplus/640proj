package main

import (
    "../delegateproto"
    "../config"
    "../coordproto"
    "flag"
    "fmt"
    "strings"
    "strconv"
    "net/rpc"
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

func parseFlight(str string) *delegateproto.FlightStruct {
    f := &delegateproto.FlightStruct{}
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

    conf, _ := config.ReadConfigFile(*configFile)
    // config.DumpConfig(conf)

    coordRPC, err := rpc.DialHTTP("tcp", conf.CoordHostPort)
    if err != nil {
        fmt.Printf("Could not connect to server %s, returning nil\n", conf.CoordHostPort)
        return
    }

    cmd := flag.Arg(0)
    for et:=0; et < *numTimes; et++ {
        switch(cmd){
        case "a":
            flightStr := flag.Arg(1)
            flight := parseFlight(flightStr)
            args := &coordproto.AddArgs{}
            var reply coordproto.AddReply
            args.Flight = *flight
            coordRPC.Call("CoordinatorRPC.AddFlight", args, &reply)
            fmt.Println("a", *args, reply)

        case "d":
            args := &coordproto.DeleteArgs{}
            var reply coordproto.DeleteReply
            args.FlightID = flag.Arg(1)
            coordRPC.Call("CoordinatorRPC.DeleteFlight", args, &reply)
            fmt.Println("d", *args, reply)

        case "q":
            args := &coordproto.QueryArgs{}
            parts := strings.Split(flag.Arg(1), ":")
            args.StartTime, _ = strconv.ParseInt(parts[0], 10, 64)
            args.EndTime, _ = strconv.ParseInt(parts[1], 10, 64)
            var reply coordproto.QueryReply
            coordRPC.Call("CoordinatorRPC.QueryFlights", args, &reply)
            fmt.Println("q", *args, reply)

        case "r":
            flightStr := flag.Arg(1)
            flight := parseFlight(flightStr)
            args := &coordproto.RescheduleArgs{}
            var reply coordproto.RescheduleReply
            args.NewFlight = *flight
            args.OldFlightID = flight.FlightID
            coordRPC.Call("CoordinatorRPC.RescheduleFlight", args, &reply)
            fmt.Println("r", *args, reply)

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
            var reply coordproto.BookReply
            coordRPC.Call("CoordinatorRPC.CancelFlights", args, &reply)
            fmt.Println("c", *args, reply)
        }
    }
    coordRPC.Close()
}
