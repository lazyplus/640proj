package delegateimpl

import (
	"time"
    "fmt"
    "net/rpc"
    "../paxosproto"
    "../config"
    "../delegateproto"
    "encoding/json"
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

func (dg *Delegate) Push(V * paxosproto.ValueStruct) *paxosproto.ReplyStruct{
    index := 0
    fmt.Println("Pushing ")
    fmt.Println(V)
    for ; ; index = (index + 1) % dg.conf.NumPeers {
//        var reply []byte
		reply := &paxosproto.ReplyStruct{}
        client, err := rpc.DialHTTP("tcp", dg.conf.PeersHostPort[index])
        if err != nil {
            // fmt.Println("Peer " + dg.conf.PeersHostPort[index] + " dead")
            continue
        }
        fmt.Println("Calling Propose to " + dg.conf.PeersHostPort[index])
        fmt.Println(V)
        err = client.Call("PaxosEngine.Propose", V, reply)
        if err != nil {
            fmt.Println(err)
            client.Close()
            continue
        }
        if reply != nil && reply.Status == paxosproto.Propose_OK {
            fmt.Println("got reply")
            fmt.Println(reply)
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
    // V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_PrepareCancelFlight
    // reply = dg.Push(V).(*delegateproto.BookReply)
    return nil
}

func (dg *Delegate) QueryFlights(args * delegateproto.QueryArgs, reply * delegateproto.QueryReply) error {
    V := paxosproto.ValueStruct{}
    // V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_QueryFlights
    // reply = dg.Push(V).(*delegateproto.QueryReply)
    return nil
}

func (dg *Delegate) PrepareBookFlight(args * delegateproto.BookArgs, reply * delegateproto.BookReply) error {
    V := paxosproto.ValueStruct{}
    // V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_PrepareBookFlight
    // reply = dg.Push(V).(*delegateproto.BookReply)
    return nil
}

func (dg *Delegate) BookDecision(args * delegateproto.DecisionArgs, reply * delegateproto.DecisionReply) error {
    V := paxosproto.ValueStruct{}
    // V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_BookDecision
    // reply = dg.Push(V).(*delegateproto.DecisionReply)
    return nil
}

func (dg *Delegate) CancelDecision(args * delegateproto.DecisionArgs, reply * delegateproto.DecisionReply) error {
    V := paxosproto.ValueStruct{}
    // V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_CancelDecision
    // reply = dg.Push(V).(*delegateproto.DecisionReply)
    return nil
}

func (dg *Delegate) DeleteFlight(args * delegateproto.DeleteArgs, reply * delegateproto.DeleteReply) error {
    V := paxosproto.ValueStruct{}
    // V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_DeleteFlight
    // reply = dg.Push(V).(*delegateproto.DeleteReply)
    return nil
}

func (dg *Delegate) RescheduleFlight(args * delegateproto.RescheduleArgs, reply * delegateproto.RescheduleReply) error {
    V := paxosproto.ValueStruct{}
    // V.Action = args
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_RescheduleFlight
    // reply = dg.Push(V).(*delegateproto.RescheduleReply)
    return nil
}

func (dg *Delegate) AddFlight(args * delegateproto.AddArgs, reply * delegateproto.AddReply) error {
    fmt.Println("AddFlight Called")
    V := &paxosproto.ValueStruct{}
    buf, _ := json.Marshal(*args)
    V.Action = make([]byte, len(buf))
    copy(V.Action, buf)
    V.Host = dg.conf.DelegateHostPort
    V.CoordSeq = args.Seqnum
    V.Type = paxosproto.C_AddFlight
    rpl := dg.Push(V)
    if rpl.Status == paxosproto.Propose_OK {
   	 err := json.Unmarshal(rpl.Reply, reply)
   	 if err != nil {
    	    fmt.Println(err)
    	}
    }
    return nil
}
