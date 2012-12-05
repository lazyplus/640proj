package coordproto

import (
    "../delegateproto"
)

// Status Codes
const (
    OK = iota
    ENOFLIGHT
    ENOTICKET
    ENOSEQ
    ENOAIRLINE
)

type airlineinfo struct {
	address string
}

type BookArgs struct {
    Flights []string
    Email string
    Count int
    Seq	int
}

type BookReply struct {
    Status int
    Seq int
}

type QueryArgs struct {
    StartTime int64
    EndTime int64
}

type QueryReply struct {
    Status int
    FlightList []delegateproto.FlightStruct
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
    NewFlight delegateproto.FlightStruct
}

type RescheduleReply struct {
    Status int
    CustomerEmails []string
}

type AddArgs struct {
    Flight delegateproto.FlightStruct
}

type AddReply struct {
    Status int
}
