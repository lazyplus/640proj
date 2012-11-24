package coordproto

// Status Codes
const (
    OK = iota
    ENOFLIGHT
    ENOTICKET
    ENOSEQ
)

type airlineinfo struct {
	address string
}

type BookArgs struct {
    Flights []string
    Email string
    Count int
}

type BookReply struct {
    Status int
}

type CancelArgs struct {
    Flights []string
    Email string
}

type CancelReply struct {
    Status int
}
