package delegateproto

// Status Codes
const (
    OK = iota
    ENOFLIGHT
    ENOTICKET
    EFLIGHTEXISTS
    ENOPREPACT
    ETEMP
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
    Seqnum int
}

type QueryReply struct {
    Status int
    FlightList []FlightStruct
    Seqnum int
}

type BookArgs struct {
    FlightID string
    Email string
    Count int
    Seqnum int
}

type BookReply struct {
    Status int
    Seqnum int
}

type DecisionArgs struct {
    Decision int
    FlightID string
    Seqnum int
}

type DecisionReply struct {
    Status int
    Seqnum int
}

type DeleteArgs struct {
    FlightID string
    Seqnum int
}

type DeleteReply struct {
    Status int
    CustomerEmails []string
    Seqnum int
}

type RescheduleArgs struct {
    OldFlightID string
    NewFlight FlightStruct
    Seqnum int
}

type RescheduleReply struct {
    Status int
    CustomerEmails []string
    Seqnum int
}

type AddArgs struct {
    Flight FlightStruct
    Seqnum int
}

type AddReply struct {
    Status int
    Seqnum int
}

type DelegateServerInterface interface {
    QueryFlights(* QueryArgs, * QueryReply) error
    PrepareBookFlight(* BookArgs, * BookReply) error
    BookDecision(* DecisionArgs, * DecisionReply) error
    PrepareCancelFlight(* BookArgs, * BookReply) error
    CancelDecision(* DecisionArgs, * DecisionReply) error
    DeleteFlight(* DeleteArgs, * DeleteReply) error
    RescheduleFlight(* RescheduleArgs, * RescheduleReply) error
    AddFlight(* AddArgs, * AddReply) error
}

type DelegateServerRPC struct {
    ds DelegateServerInterface
}

func NewDelegateServerRPC(ds DelegateServerInterface) *DelegateServerRPC {
    return &DelegateServerRPC{ds}
}

func (dsrpc *DelegateServerRPC) QueryFlights(args * QueryArgs, reply * QueryReply) error {
    return dsrpc. ds.QueryFlights(args, reply)
}

func (dsrpc *DelegateServerRPC) PrepareBookFlight(args * BookArgs, reply * BookReply) error {
    return dsrpc. ds.PrepareBookFlight(args, reply)
}

func (dsrpc *DelegateServerRPC) BookDecision(args * DecisionArgs, reply * DecisionReply) error {
    return dsrpc. ds.BookDecision(args, reply)
}

func (dsrpc *DelegateServerRPC) PrepareCancelFlight(args * BookArgs, reply * BookReply) error {
    return dsrpc. ds.PrepareCancelFlight(args, reply)
}

func (dsrpc *DelegateServerRPC) CancelDecision(args * DecisionArgs, reply * DecisionReply) error {
    return dsrpc. ds.CancelDecision(args, reply)
}

func (dsrpc *DelegateServerRPC) DeleteFlight(args * DeleteArgs, reply * DeleteReply) error {
    return dsrpc. ds.DeleteFlight(args, reply)
}

func (dsrpc *DelegateServerRPC) RescheduleFlight(args * RescheduleArgs, reply * RescheduleReply) error {
    return dsrpc. ds.RescheduleFlight(args, reply)
}

func (dsrpc *DelegateServerRPC) AddFlight(args * AddArgs, reply * AddReply) error {
    return dsrpc. ds.AddFlight(args, reply)
}
