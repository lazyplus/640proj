package airlineserver

import (
    "../ultility"
    "../delegateproto"
    "sync"
    "fmt"
)

type FlightInfo struct {
    flight *delegateproto.FlightStruct
    preparedAction *delegateproto.BookArgs
    customers map[string] int
    mutex ultility.PLock
    deleted bool
}

type AirlineServer struct {
    flightListLock sync.Mutex
    flightList map[string] *FlightInfo
}

func NewAirlineServer () *AirlineServer {
    as := &AirlineServer{}
    as.flightList = make(map[string] *FlightInfo)
    return as
}

func (as *AirlineServer) Progress(V ValueStruct) interface{}, error {
    var reply interface{}
    switch (V.Type) {
    case c_QueryFlights:
        err := as.QueryFlights(V.Action.(delegateproto.QueryArgs), reply)
        if err != nil {
            return nil, err
        }
    case c_PrepareBookFlight:
        err := as.PrepareBookFlight(V.Action.(delegateproto.BookArgs), reply)
        if err != nil {
            return nil, err
        }
    case c_PrepareCancelFlight:
        err := as.PrepareCancelFlight(V.Action.(delegateproto.BookArgs), reply)
        if err != nil {
            return nil, err
        }
    case c_BookDecision:
        err := as.BookDecision(V.Action.(delegateproto.BookDecision), reply)
        if err != nil {
            return nil, err
        }
    case c_CancelDecision:
        err := as.CancelDecision(V.Action.(delegateproto.BookDecision), reply)
        if err != nil {
            return nil, err
        }
    case c_DeleteFlight:
        err := as.DeleteFlight(V.Action.(delegateproto.DeleteArgs), reply)
        if err != nil {
            return nil, err
        }
    case c_RescheduleFlight:
        err := as.RescheduleFlight(V.Action.(delegateproto.RescheduleArgs), reply)
        if err != nil {
            return nil, err
        }
    case c_AddFlight:
        err := as.AddFlight(V.Action.(delegateproto.AddArgs), reply)
        if err != nil {
            return nil, err
        }
    }
    return reply, nil
}

func (as *AirlineServer) getFlight(id string) *FlightInfo {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    return as.flightList[id]
}

func (as *AirlineServer) QueryFlights(args *delegateproto.QueryArgs, reply *delegateproto.QueryReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    reply.FlightList = make([]delegateproto.FlightStruct, 0)
    for _, value := range as.flightList {
        if value.flight.DepartureTime <= args.EndTime && value.flight.DepartureTime >= args.StartTime {
            reply.FlightList = append(reply.FlightList, *value.flight)
        }
    }

    reply.Status = delegateproto.OK
    return nil
}

func (as *AirlineServer) PrepareBookFlight(args *delegateproto.BookArgs, reply *delegateproto.BookReply) error {
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return nil
    }

    getLock := flight.mutex.TryLock()

    if !getLock {
        return errors.New("Cannot get lock")
    }

    flight.preparedAction = args

    if flight.deleted {
        reply.Status = delegateproto.ENOFLIGHT
        return nil
    }

    if flight.flight.AvailableTickets < args.Count {
        fmt.Printf("AvailableTickets: %d\n", flight.flight.AvailableTickets)
        reply.Status = delegateproto.ENOTICKET
        return nil
    }

    reply.Status = delegateproto.OK
    return nil
}

func (as *AirlineServer) BookDecision(args *delegateproto.DecisionArgs, reply *delegateproto.DecisionReply) error {
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return nil
    }

    defer flight.mutex.Unlock()

    act := flight.preparedAction
    flight.preparedAction = nil

    if args.Decision == delegateproto.COMMIT {
        if act == nil {
            reply.Status = delegateproto.ENOPREPACT
            return nil
        }
        flight.customers[act.Email] += act.Count
        flight.flight.AvailableTickets -= act.Count
    }

    reply.Status = delegateproto.OK
    return nil
}

func (as *AirlineServer) PrepareCancelFlight(args *delegateproto.BookArgs, reply *delegateproto.BookReply) error {
    fmt.Println("Called PrepareCancelFlight " + args.FlightID)
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return nil
    }

    getLock := flight.mutex.TryLock()

    if !getLock {
        return errors.New("Cannot get lock")
    }

    flight.preparedAction = args

    if flight.deleted {
        reply.Status = delegateproto.ENOFLIGHT
        return nil
    }

    if flight.customers[args.Email] < args.Count {
        reply.Status = delegateproto.ENOTICKET
        return nil
    }

    reply.Status = delegateproto.OK
    return nil
}

func (as *AirlineServer) CancelDecision(args *delegateproto.DecisionArgs, reply *delegateproto.DecisionReply) error {
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return nil
    }

    defer flight.mutex.Unlock()

    act := flight.preparedAction
    flight.preparedAction = nil

    if args.Decision == delegateproto.COMMIT {
        if act == nil {
            reply.Status = delegateproto.ENOPREPACT
            return nil
        }
        flight.customers[act.Email] -= act.Count
        if flight.customers[act.Email] == 0 {
            delete(flight.customers, act.Email)
        }
        flight.flight.AvailableTickets += act.Count
    }
    
    reply.Status = delegateproto.OK
    return nil
}

func (as *AirlineServer) DeleteFlight(args *delegateproto.DeleteArgs, reply *delegateproto.DeleteReply) error {
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return nil
    }

    getLock := flight.mutex.TryLock()

    if !getLock {
        return errors.New("Cannot get lock")
    }

    defer flight.mutex.Unlock()

    as.flightListLock.Lock()
    delete(as.flightList, args.FlightID)
    as.flightListLock.Unlock()

    reply.CustomerEmails = make([]string, 0)
    for key, _ := range flight.customers {
        reply.CustomerEmails = append(reply.CustomerEmails, key)
    }
    flight.deleted = true
    reply.Status = delegateproto.OK
    return nil
}

func (as *AirlineServer) RescheduleFlight(args *delegateproto.RescheduleArgs, reply *delegateproto.RescheduleReply) error {
    flight := as.getFlight(args.OldFlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return nil
    }

    getLock := flight.mutex.TryLock()

    if !getLock {
        return errors.New("Cannot get lock")
    }

    reply.CustomerEmails = make([]string, 0)
    for key, _ := range flight.customers {
        reply.CustomerEmails = append(reply.CustomerEmails, key)
    }
    flight.flight = &args.NewFlight
    reply.Status = delegateproto.OK

    flight.mutex.Unlock()
    return nil
}

func (as *AirlineServer) AddFlight(args *delegateproto.AddArgs, reply *delegateproto.AddReply) error {
    flight := as.getFlight(args.Flight.FlightID)

    if flight != nil {
        reply.Status = delegateproto.EFLIGHTEXISTS
        return nil
    }

    flightInfo := &FlightInfo{}
    flightInfo.flight = &args.Flight
    flightInfo.deleted = false
    flightInfo.customers = make(map[string] int)
    flightInfo.preparedAction = nil

    as.flightListLock.Lock()

    as.flightList[args.Flight.FlightID] = flightInfo
    reply.Status = delegateproto.OK

    as.flightListLock.Unlock()
    return nil
}
