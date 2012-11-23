package airlinerpc

import "../airlineproto"

type AirlineServerInterface interface {
    QueryFlights(*airlineproto.QueryArgs, *airlineproto.QueryReply) error
    PrepareBookFlight(*airlineproto.BookArgs, *airlineproto.BookReply) error
    BookDecision(*airlineproto.DecisionArgs, *airlineproto.DecisionReply) error
    PrepareCancelFlight(*airlineproto.CancelArgs, *airlineproto.CancelReply) error
    CancelDecision(*airlineproto.DecisionArgs, *airlineproto.DecisionReply) error
    DeleteFlight(*airlineproto.DeleteArgs, *airlineproto.DeleteReply) error
    RescheduleFlight(*airlineproto.RescheduleArgs, *airlineproto.RescheduleReply) error
    AddFlight(*airlineproto.AddArgs, *airlineproto.AddReply) error
}

type AirlineServerRPC struct {
    as AirlineServerInterface
}

func NewAirlineServerRPC(as AirlineServerInterface) *AirlineServerRPC {
    return &AirlineServerRPC{as}
}

func (asrpc *AirlineServerRPC) QueryFlights(args *airlineproto.QueryArgs, reply *airlineproto.QueryReply) error {
    return asrpc.as.QueryFlights(args, reply)
}

func (asrpc *AirlineServerRPC) PrepareBookFlight(args *airlineproto.BookArgs, reply *airlineproto.BookReply) error {
    return asrpc.as.BookFlight(args, reply)
}

func (asrpc *AirlineServerRPC) BookDecision(*airlineproto.DecisionArgs, *airlineproto.DecisionReply) error {
    return asrpc.as.BookDecision(args, reply)
}

func (asrpc *AirlineServerRPC) PrepareCancelFlight(args *airlineproto.CancelArgs, reply *airlineproto.CancelReply) error {
    return asrpc.as.CancelFlight(args, reply)
}

func (asrpc *AirlineServerRPC) CancelDecision(*airlineproto.DecisionArgs, *airlineproto.DecisionReply) error {
    return asrpc.as.CancelDecision(args, reply)
}

func (asrpc *AirlineServerRPC) DeleteFlight(args *airlineproto.DeleteArgs, reply *airlineproto.DeleteReply) error {
    return asrpc.as.DeleteFlight(args, reply)
}

func (asrpc *AirlineServerRPC) RescheduleFlight(args *airlineproto.RescheduleArgs, reply *airlineproto.RescheduleReply) error {
    return asrpc.as.RescheduleFlight(args, reply)
}

func (asrpc *AirlineServerRPC) AddFlight(args *airlineproto.AddArgs, reply *airlineproto.AddReply) error {
    return asrpc.as.AddFlight(args, reply)
}
