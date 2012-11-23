package airlineimpl

import (
    "../airlineproto"
    "sync"
)

type AirlineServer struct {
    flightListLock sync.Mutex
    flightList map[string] *airlineproto.FlightStruct
    customers map[string] map[string] int
    prepareAction map[int] *airlineproto.BookArgs
    commitSeq int
}

func NewAirlineServer () *AirlineServer {
    as := &AirlineServer{}
    as.commitSeq = 0
    as.flightList = make(map[string] *airlineproto.FlightStruct)
    as.customers = make(map[string] map[string] int)
    as.prepareAction = make(map[int] *airlineproto.BookArgs)
    return as
}

func (as *AirlineServer) QueryFlights(args *airlineproto.QueryArgs, reply *airlineproto.QueryReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    reply.FlightList = make([]airlineproto.FlightStruct, 0)
    for _, value := range as.flightList {
        if value.DepartureTime <= args.EndTime && value.DepartureTime >= args.StartTime {
            reply.FlightList = append(reply.FlightList, *value)
        }
    }

    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) PrepareBookFlight(args *airlineproto.BookArgs, reply *airlineproto.BookReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    flight := as.flightList[args.FlightID]

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }
    if flight.AvailableTickets < args.Count {
        reply.Status = airlineproto.ENOTICKET
        return nil
    }

    reply.Status = airlineproto.OK
    reply.Seq = as.commitSeq
    as.commitSeq ++
    as.prepareAction[reply.Seq] = args
    return nil
}

func (as *AirlineServer) BookDecision(args *airlineproto.DecisionArgs, reply *airlineproto.DecisionReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    if as.prepareAction[args.Seq] == nil {
        reply.Status = airlineproto.ENOSEQ
        return nil
    }

    act := as.prepareAction[args.Seq]
    as.customers[act.FlightID][act.Email] += act.Count
    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) PrepareCancelFlight(args *airlineproto.BookArgs, reply *airlineproto.BookReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    flight := as.flightList[args.FlightID]

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    reply.Status = airlineproto.OK
    reply.Seq = as.commitSeq
    as.commitSeq ++
    as.prepareAction[reply.Seq] = args
    return nil
}

func (as *AirlineServer) CancelDecision(args *airlineproto.DecisionArgs, reply *airlineproto.DecisionReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    if as.prepareAction[args.Seq] == nil {
        reply.Status = airlineproto.ENOSEQ
        return nil
    }

    act := as.prepareAction[args.Seq]
    as.customers[act.FlightID][act.Email] -= act.Count
    if as.customers[act.FlightID][act.Email] == 0 {
        delete(as.customers[act.FlightID], act.Email)
    }
    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) DeleteFlight(args *airlineproto.DeleteArgs, reply *airlineproto.DeleteReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    flight := as.flightList[args.FlightID]

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    reply.CustomerEmails = make([]string, 0)
    for key, _ := range as.customers[args.FlightID] {
        reply.CustomerEmails = append(reply.CustomerEmails, key)
    }
    delete(as.flightList, args.FlightID)
    delete(as.customers, args.FlightID)
    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) RescheduleFlight(args *airlineproto.RescheduleArgs, reply *airlineproto.RescheduleReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    flight := as.flightList[args.OldFlightID]

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    reply.CustomerEmails = make([]string, 0)
    for key, _ := range as.customers[args.OldFlightID] {
        reply.CustomerEmails = append(reply.CustomerEmails, key)
    }
    as.flightList[args.OldFlightID] = &args.NewFlight
    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) AddFlight(args *airlineproto.AddArgs, reply *airlineproto.AddReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    as.flightList[args.Flight.FlightID] = &args.Flight
    as.customers[args.Flight.FlightID] = make(map[string] int)
    reply.Status = airlineproto.OK
    return nil
}
