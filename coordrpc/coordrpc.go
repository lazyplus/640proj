package coordrpc

import "../coordproto"

type CoordinatorInterface interface {
    BookFlights(*coordproto.BookArgs, *coordproto.BookReply) error
    CancelFlights(*coordproto.BookArgs, *coordproto.BookReply) error
    QueryFlights(* coordproto.QueryArgs, * coordproto.QueryReply) error
    DeleteFlight(* coordproto.DeleteArgs, * coordproto.DeleteReply) error
    RescheduleFlight(* coordproto.RescheduleArgs, * coordproto.RescheduleReply) error
    AddFlight(* coordproto.AddArgs, * coordproto.AddReply) error
}

type CoordinatorRPC struct {
    cd CoordinatorInterface
}

func NewCoordinatorRPC(cd CoordinatorInterface) *CoordinatorRPC {
    return &CoordinatorRPC{cd}
}

func (cdrpc *CoordinatorRPC) BookFlights(args *coordproto.BookArgs, reply *coordproto.BookReply) error {
    return cdrpc.cd.BookFlights(args, reply)
}

func (cdrpc *CoordinatorRPC) CancelFlights(args *coordproto.BookArgs, reply *coordproto.BookReply) error {
    return cdrpc.cd.CancelFlights(args, reply)
}

func (cdrpc *CoordinatorRPC) QueryFlights(args * coordproto.QueryArgs, reply * coordproto.QueryReply) error {
    return cdrpc.cd.QueryFlights(args, reply)
}

func (cdrpc *CoordinatorRPC) DeleteFlight(args * coordproto.DeleteArgs, reply * coordproto.DeleteReply) error {
    return cdrpc.cd.DeleteFlight(args, reply)
}

func (cdrpc *CoordinatorRPC) RescheduleFlight(args * coordproto.RescheduleArgs, reply * coordproto.RescheduleReply) error {
    return cdrpc.cd.RescheduleFlight(args, reply)
}

func (cdrpc *CoordinatorRPC) AddFlight(args * coordproto.AddArgs, reply * coordproto.AddReply) error {
    return cdrpc.cd.AddFlight(args, reply)
}
