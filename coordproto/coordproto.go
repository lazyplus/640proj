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
    Seq	int
}

type BookReply struct {
    Status int
    Seq int
}
