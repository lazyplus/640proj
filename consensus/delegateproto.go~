package consensus

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
