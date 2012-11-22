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

type DeleteArgs struct {
    FlightID string
}

type DeleteReply struct {
    Status int
    CustomerEmails []string
}

type RescheduleArgs struct {
    OldFlightID string
    NewFlight FlightStruct
}

type RescheduleReply struct {
    Status int
    CustomerEmails []string
}

type AddArgs struct {
    Flight FlightStruct
}

type AddReply struct {
    Status int
}
