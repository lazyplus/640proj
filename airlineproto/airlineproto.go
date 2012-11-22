package airlineproto

// Status Codes
const (
    OK = iota
    ENOFLIGHT
    ENOTICKET
    EFLIGHTEXISTS
)

type FlightStruct struct {
    FlightID string
    DepartureTime int64
    ArrivalTime int64
    DeparturePort string
    ArrivalPort string
    AvailableTickets int
}

type QueryArgs struct {
    StartTime int64
    EndTime int64
}

type QueryReply struct {
    Status int
    FlightList []FlightStruct
}

type BookArgs struct {
    FlightID string
    Email string
}

type BookReply struct {
    Status int
}

type CancelArgs struct {
    FlightID string
    Email string
}

type CancelReply struct {
    Status int
}

type DeleteFlightArgs struct {
    FlightID string
}

type DeleteFlightReply struct {
    Status int
    CustomerEmails []string
}

type RescheduleFlightArgs struct {
    OldFlightID string
    NewFlight FlightStruct
}

type RescheduleFlightReply struct {
    Status int
    CustomerEmails []string
}

type AddFlightArgs struct {
    Flight FlightStruct
}

type AddFlightReply struct {
    Status int
}
