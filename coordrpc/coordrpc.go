package coordrpc

import "../coordproto"

type CoordinatorInterface interface {
    BookFlights(*coordproto.BookArgs, *coordproto.BookReply) error
    CancelFlights(*coordproto.BookArgs, *coordproto.BookReply) error
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
