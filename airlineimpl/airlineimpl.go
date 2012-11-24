package airlineimpl

import (
    "../airlineproto"
    "sync"
    "fmt"
)

type FlightInfo struct {
    flight *airlineproto.FlightStruct
    preparedAction *airlineproto.BookArgs
    customers map[string] int
    mutex sync.Mutex
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

func (as *AirlineServer) getFlight(id string) *FlightInfo {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    return as.flightList[id]
}

func (as *AirlineServer) QueryFlights(args *airlineproto.QueryArgs, reply *airlineproto.QueryReply) error {
    as.flightListLock.Lock()
    defer as.flightListLock.Unlock()

    reply.FlightList = make([]airlineproto.FlightStruct, 0)
    for _, value := range as.flightList {
        if value.flight.DepartureTime <= args.EndTime && value.flight.DepartureTime >= args.StartTime {
            reply.FlightList = append(reply.FlightList, *value.flight)
        }
    }

    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) PrepareBookFlight(args *airlineproto.BookArgs, reply *airlineproto.BookReply) error {
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    flight.mutex.Lock()
    flight.preparedAction = args

    if flight.deleted {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    if flight.flight.AvailableTickets < args.Count {
        fmt.Printf("AvailableTickets: %d\n", flight.flight.AvailableTickets)
        reply.Status = airlineproto.ENOTICKET
        return nil
    }

    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) BookDecision(args *airlineproto.DecisionArgs, reply *airlineproto.DecisionReply) error {
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    defer flight.mutex.Unlock()

    act := flight.preparedAction
    flight.preparedAction = nil

    if args.Decision == airlineproto.COMMIT {
        if act == nil {
            reply.Status = airlineproto.ENOPREPACT
            return nil
        }
        flight.customers[act.Email] += act.Count
        flight.flight.AvailableTickets -= act.Count
    }

    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) PrepareCancelFlight(args *airlineproto.BookArgs, reply *airlineproto.BookReply) error {
    fmt.Println("Called PrepareCancelFlight " + args.FlightID)
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    flight.mutex.Lock()
    flight.preparedAction = args

    if flight.deleted {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    if flight.customers[args.Email] < args.Count {
        reply.Status = airlineproto.ENOTICKET
        return nil
    }

    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) CancelDecision(args *airlineproto.DecisionArgs, reply *airlineproto.DecisionReply) error {
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    defer flight.mutex.Unlock()

    act := flight.preparedAction
    flight.preparedAction = nil

    if args.Decision == airlineproto.COMMIT {
        if act == nil {
            reply.Status = airlineproto.ENOPREPACT
            return nil
        }
        flight.customers[act.Email] -= act.Count
        if flight.customers[act.Email] == 0 {
            delete(flight.customers, act.Email)
        }
        flight.flight.AvailableTickets += act.Count
    }
    
    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) DeleteFlight(args *airlineproto.DeleteArgs, reply *airlineproto.DeleteReply) error {
    flight := as.getFlight(args.FlightID)

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    flight.mutex.Lock()
    defer flight.mutex.Unlock()

    as.flightListLock.Lock()
    delete(as.flightList, args.FlightID)
    as.flightListLock.Unlock()

    reply.CustomerEmails = make([]string, 0)
    for key, _ := range flight.customers {
        reply.CustomerEmails = append(reply.CustomerEmails, key)
    }
    flight.deleted = true
    reply.Status = airlineproto.OK
    return nil
}

func (as *AirlineServer) RescheduleFlight(args *airlineproto.RescheduleArgs, reply *airlineproto.RescheduleReply) error {
    flight := as.getFlight(args.OldFlightID)

    if flight == nil {
        reply.Status = airlineproto.ENOFLIGHT
        return nil
    }

    flight.mutex.Lock()

    reply.CustomerEmails = make([]string, 0)
    for key, _ := range flight.customers {
        reply.CustomerEmails = append(reply.CustomerEmails, key)
    }
    flight.flight = &args.NewFlight
    reply.Status = airlineproto.OK

    flight.mutex.Unlock()
    return nil
}

func (as *AirlineServer) AddFlight(args *airlineproto.AddArgs, reply *airlineproto.AddReply) error {
    flight := as.getFlight(args.Flight.FlightID)

    if flight != nil {
        reply.Status = airlineproto.EFLIGHTEXISTS
        return nil
    }

    flightInfo := &FlightInfo{}
    flightInfo.flight = &args.Flight
    flightInfo.deleted = false
    flightInfo.customers = make(map[string] int)
    flightInfo.preparedAction = nil

    as.flightListLock.Lock()

    as.flightList[args.Flight.FlightID] = flightInfo
    reply.Status = airlineproto.OK

    as.flightListLock.Unlock()
    return nil
}
