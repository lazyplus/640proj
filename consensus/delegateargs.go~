package consensus

// Status Codes
const (
    OK = iota
    ENOFLIGHT
    ENOTICKET
    EFLIGHTEXISTS
    ENOPREPACT
)

const (
    COMMIT = iota
    ABORT
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

type QueryArgs struct {
    StartTime int64
    EndTime int64
    seqnum int
}

type QueryReply struct {
    Status int
    FlightList []FlightStruct
    seqnum int
}

type BookArgs struct {
    FlightID string
    Email string
    Count int
    seqnum int
}

type BookReply struct {
    Status int
    seqnum int
}

type DecisionArgs struct {
    Decision int
    FlightID string
	seqnum int
}

type DecisionReply struct {
    Status int
    seqnum int
}

type DeleteArgs struct {
    FlightID string
    seqnum int
}

type DeleteReply struct {
    Status int
    CustomerEmails []string
    seqnum int
}

type RescheduleArgs struct {
    OldFlightID string
    NewFlight FlightStruct
    seqnum int
}

type RescheduleReply struct {
    Status int
    CustomerEmails []string
    seqnum int
}

type AddArgs struct {
    Flight FlightStruct
    seqnum int
}

type AddReply struct {
    Status int
    seqnum int
}
