package delegateimpl

import (
	"time"
    "net/rpc"
    "../paxosproto"
    "../config"
    "../delegateproto"
)

type Delegate struct {
    conf *config.AirlineConfig
}

func NewDelegate(path string, airline_name string) *Delegate {
    dg := &Delegate{}
	cg, _ := config.ReadConfigFile(path)
    ac, found := cg.Airlines[airline_name]
	if !found {
		return nil
	}
    dg.conf = ac
	return dg
}

func (dg *Delegate) Push(V paxosproto.ValueStruct) interface{} {
    index := 0

    for ; ; index = (index + 1) % dg.conf.NumPeers {
        var reply interface{}
        client, err := rpc.DialHTTP("tcp", dg.conf.PeersHostPort[index])
        if err != nil {
            continue
        }
        err = client.Call("AirlineServerRPC.Propose", V, reply)
        if err != nil {
            client.Close()
            continue
        }
        if reply != nil {
            client.Close()
            return reply
        }
        time.Sleep(time.Second)
    }
    return nil
}

func (dg *Delegate) PrepareCancelFlight(args * delegateproto.BookArgs, reply * delegateproto.BookReply) error {
    // make up Value
    V := paxosproto.ValueStruct{}
    V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_PrepareCancelFlight
    reply = dg.Push(V).(*delegateproto.BookReply)
    return nil
}

func (dg *Delegate) QueryFlights(args * delegateproto.QueryArgs, reply * delegateproto.QueryReply) error {
    V := paxosproto.ValueStruct{}
    V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_QueryFlights
    reply = dg.Push(V).(*delegateproto.QueryReply)
    return nil
}

func (dg *Delegate) PrepareBookFlight(args * delegateproto.BookArgs, reply * delegateproto.BookReply) error {
    V := paxosproto.ValueStruct{}
    V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_PrepareBookFlight
    reply = dg.Push(V).(*delegateproto.BookReply)
    return nil
}

func (dg *Delegate) BookDecision(args * delegateproto.DecisionArgs, reply * delegateproto.DecisionReply) error {
    V := paxosproto.ValueStruct{}
    V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_BookDecision
    reply = dg.Push(V).(*delegateproto.DecisionReply)
    return nil
}

func (dg *Delegate) CancelDecision(args * delegateproto.DecisionArgs, reply * delegateproto.DecisionReply) error {
    V := paxosproto.ValueStruct{}
    V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_CancelDecision
    reply = dg.Push(V).(*delegateproto.DecisionReply)
    return nil
}

func (dg *Delegate) DeleteFlight(args * delegateproto.DeleteArgs, reply * delegateproto.DeleteReply) error {
    V := paxosproto.ValueStruct{}
    V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_DeleteFlight
    reply = dg.Push(V).(*delegateproto.DeleteReply)
    return nil
}

func (dg *Delegate) RescheduleFlight(args * delegateproto.RescheduleArgs, reply * delegateproto.RescheduleReply) error {
    V := paxosproto.ValueStruct{}
    V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_RescheduleFlight
    reply = dg.Push(V).(*delegateproto.RescheduleReply)
    return nil
}

func (dg *Delegate) AddFlight(args * delegateproto.AddArgs, reply * delegateproto.AddReply) error {
    V := paxosproto.ValueStruct{}
    V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_AddFlight
    reply = dg.Push(V).(*delegateproto.AddReply)
    return nil
}
