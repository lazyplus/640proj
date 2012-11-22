package coordproto

// Status Codes
const (
    OK = iota
    ENOFLIGHT
    ENOTICKET
)

type BookArgs struct {
    Flights []string
    Email string
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
