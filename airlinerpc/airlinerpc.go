package airlinerpc

import "../airlineproto"

type AirlineServerInterface interface {
    QueryFlights(*airlineproto.QueryArgs, *airlineproto.QueryReply) error
    PrepareBookFlight(*airlineproto.BookArgs, *airlineproto.BookReply) error
    BookDecision(*airlineproto.DecisionArgs, *airlineproto.DecisionReply) error
    PrepareCancelFlight(*airlineproto.BookArgs, *airlineproto.BookReply) error
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
    return asrpc.as.PrepareBookFlight(args, reply)
}

func (asrpc *AirlineServerRPC) BookDecision(args *airlineproto.DecisionArgs, reply *airlineproto.DecisionReply) error {
    return asrpc.as.BookDecision(args, reply)
}

func (asrpc *AirlineServerRPC) PrepareCancelFlight(args *airlineproto.BookArgs, reply *airlineproto.BookReply) error {
    return asrpc.as.PrepareCancelFlight(args, reply)
}

func (asrpc *AirlineServerRPC) CancelDecision(args *airlineproto.DecisionArgs, reply *airlineproto.DecisionReply) error {
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
