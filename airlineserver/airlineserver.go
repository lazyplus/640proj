package airlineserver

import (
    "../ultility"
    "../delegateproto"
    "../paxosproto"
    "sync"
    "fmt"
    "errors"
    "encoding/json"
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

func (as *AirlineServer) Progress(V *paxosproto.ValueStruct) (interface{}, error) {
    var reply interface{}
    var err error

    switch (V.Type) {
    case paxosproto.C_QueryFlights:
        // reply, err = as.QueryFlights(V.Action.(delegateproto.QueryArgs))
        if err != nil {
            return nil, err
        }
    case paxosproto.C_PrepareBookFlight:
        // reply, err = as.PrepareBookFlight(V.Action.(delegateproto.BookArgs))
        if err != nil {
            return nil, err
        }
    case paxosproto.C_PrepareCancelFlight:
        // reply, err = as.PrepareCancelFlight(V.Action.(delegateproto.BookArgs))
        if err != nil {
            return nil, err
        }
    case paxosproto.C_BookDecision:
        // reply, err = as.BookDecision(V.Action.(delegateproto.DecisionArgs))
        if err != nil {
            return nil, err
        }
    case paxosproto.C_CancelDecision:
        // reply, err = as.CancelDecision(V.Action.(delegateproto.DecisionArgs))
        if err != nil {
            return nil, err
        }
    case paxosproto.C_DeleteFlight:
        // reply, err = as.DeleteFlight(V.Action.(delegateproto.DeleteArgs))
        if err != nil {
            return nil, err
        }
    case paxosproto.C_RescheduleFlight:
        // reply, err = as.RescheduleFlight(V.Action.(delegateproto.RescheduleArgs))
        if err != nil {
            return nil, err
        }
    case paxosproto.C_AddFlight:
        data := &delegateproto.AddArgs{}
        err := json.Unmarshal(V.Action, data)
        if err != nil {
            fmt.Println(err)
            return nil, err
        }

        reply, err = as.AddFlight(data)
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

func (as *AirlineServer) QueryFlights(args delegateproto.QueryArgs) (*delegateproto.QueryReply, error) {
    reply := &delegateproto.QueryReply{}
    reply.Seqnum = args.Seqnum

    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    reply.FlightList = make([]delegateproto.FlightStruct, 0)
    for _, value := range as.flightList {
        if value.flight.DepartureTime <= args.EndTime && value.flight.DepartureTime >= args.StartTime {
            reply.FlightList = append(reply.FlightList, *value.flight)
        }
    }

    reply.Status = delegateproto.OK
    return reply, nil
}

func (as *AirlineServer) PrepareBookFlight(args delegateproto.BookArgs) (*delegateproto.BookReply, error) {
    reply := &delegateproto.BookReply{}
    reply.Seqnum = args.Seqnum

    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return reply, nil
    }

    getLock := flight.mutex.TryLock()

    if !getLock {
        return nil, errors.New("Cannot get lock")
    }

    /// WATCHOUT
    flight.preparedAction = &args

    if flight.deleted {
        reply.Status = delegateproto.ENOFLIGHT
        return reply, nil
    }

    if flight.flight.AvailableTickets < args.Count {
        fmt.Printf("AvailableTickets: %d\n", flight.flight.AvailableTickets)
        reply.Status = delegateproto.ENOTICKET
        return reply, nil
    }

    reply.Status = delegateproto.OK
    return reply, nil
}

func (as *AirlineServer) BookDecision(args delegateproto.DecisionArgs) (*delegateproto.DecisionReply, error) {
    reply := &delegateproto.DecisionReply{}
    reply.Seqnum = args.Seqnum

    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return reply, nil
    }

    defer flight.mutex.Unlock()

    act := flight.preparedAction
    flight.preparedAction = nil

    if args.Decision == delegateproto.COMMIT {
        if act == nil {
            reply.Status = delegateproto.ENOPREPACT
            return reply, nil
        }
        flight.customers[act.Email] += act.Count
        flight.flight.AvailableTickets -= act.Count
    }

    reply.Status = delegateproto.OK
    return reply, nil
}

func (as *AirlineServer) PrepareCancelFlight(args delegateproto.BookArgs) (*delegateproto.BookReply, error) {
    reply := &delegateproto.BookReply{}
    reply.Seqnum = args.Seqnum

    fmt.Println("Called PrepareCancelFlight " + args.FlightID)
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return reply, nil
    }

    getLock := flight.mutex.TryLock()

    if !getLock {
        return nil, errors.New("Cannot get lock")
    }

    /// WATCHOUT
    flight.preparedAction = &args

    if flight.deleted {
        reply.Status = delegateproto.ENOFLIGHT
        return reply, nil
    }

    if flight.customers[args.Email] < args.Count {
        reply.Status = delegateproto.ENOTICKET
        return reply, nil
    }

    reply.Status = delegateproto.OK
    return reply, nil
}

func (as *AirlineServer) CancelDecision(args delegateproto.DecisionArgs) (*delegateproto.DecisionReply, error) {
    reply := &delegateproto.DecisionReply{}
    reply.Seqnum = args.Seqnum

    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return reply, nil
    }

    defer flight.mutex.Unlock()

    act := flight.preparedAction
    flight.preparedAction = nil

    if args.Decision == delegateproto.COMMIT {
        if act == nil {
            reply.Status = delegateproto.ENOPREPACT
            return reply, nil
        }
        flight.customers[act.Email] -= act.Count
        if flight.customers[act.Email] == 0 {
            delete(flight.customers, act.Email)
        }
        flight.flight.AvailableTickets += act.Count
    }
    
    reply.Status = delegateproto.OK
    return reply, nil
}

func (as *AirlineServer) DeleteFlight(args delegateproto.DeleteArgs) (*delegateproto.DeleteReply, error) {
    reply := &delegateproto.DeleteReply{}
    reply.Seqnum = args.Seqnum

    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return reply, nil
    }

    getLock := flight.mutex.TryLock()

    if !getLock {
        return nil, errors.New("Cannot get lock")
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
    return reply, nil
}

func (as *AirlineServer) RescheduleFlight(args delegateproto.RescheduleArgs) (*delegateproto.RescheduleReply, error) {
    reply := &delegateproto.RescheduleReply{}
    reply.Seqnum = args.Seqnum

    flight := as.getFlight(args.OldFlightID)

    if flight == nil {
        reply.Status = delegateproto.ENOFLIGHT
        return reply, nil
    }

    getLock := flight.mutex.TryLock()

    if !getLock {
        return nil, errors.New("Cannot get lock")
    }

    reply.CustomerEmails = make([]string, 0)
    for key, _ := range flight.customers {
        reply.CustomerEmails = append(reply.CustomerEmails, key)
    }
    flight.flight = &args.NewFlight
    reply.Status = delegateproto.OK

    flight.mutex.Unlock()
    return reply, nil
}

func (as *AirlineServer) AddFlight(args *delegateproto.AddArgs) (*delegateproto.AddReply, error) {
    reply := &delegateproto.AddReply{}
    reply.Seqnum = args.Seqnum

    flight := as.getFlight(args.Flight.FlightID)

    if flight != nil {
        reply.Status = delegateproto.EFLIGHTEXISTS
        return reply, nil
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
    return reply, nil
}
