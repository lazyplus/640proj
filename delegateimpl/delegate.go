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
    for index := 0; ; index = (index + 1) % dg.conf.NumPeers {
		reply := &paxosproto.ReplyStruct{}
        client, err := rpc.DialHTTP("tcp", dg.conf.PeersHostPort[index])
        if err != nil {
            continue
        }
        // fmt.Println("Calling Propose to " + dg.conf.PeersHostPort[index])
        err = client.Call("PaxosEngine.Propose", V, reply)
        if err != nil {
            fmt.Println(err)
            client.Close()
            continue
        }
        if reply.Status == paxosproto.Propose_OK {
            // fmt.Println("got reply")
            client.Close()
            return reply
        }
        if reply.Status == paxosproto.Propose_RETRY {
            // fmt.Println("got Retry")
            index = index + dg.conf.NumPeers - 1
            client.Close()
            continue
        }
        time.Sleep(time.Second)
    }
    return nil
}

func (dg *Delegate) PrepareCancelFlight(args * delegateproto.BookArgs, reply * delegateproto.BookReply) error {
    r := dg.doPush(*args, paxosproto.C_PrepareCancelFlight, args.Seqnum, *reply)
    err := json.Unmarshal(r, reply)
    if err != nil {
        fmt.Println(err)
    }
    return nil
}

func (dg *Delegate) PrepareBookFlight(args * delegateproto.BookArgs, reply * delegateproto.BookReply) error {
    r := dg.doPush(*args, paxosproto.C_PrepareBookFlight, args.Seqnum, *reply)
    err := json.Unmarshal(r, reply)
    if err != nil {
        fmt.Println(err)
    }
    return nil
}

func (dg *Delegate) BookDecision(args * delegateproto.DecisionArgs, reply * delegateproto.DecisionReply) error {
    r := dg.doPush(*args, paxosproto.C_BookDecision, args.Seqnum, *reply)
    err := json.Unmarshal(r, reply)
    if err != nil {
        fmt.Println(err)
    }
    return nil
}

func (dg *Delegate) CancelDecision(args * delegateproto.DecisionArgs, reply * delegateproto.DecisionReply) error {
    r := dg.doPush(*args, paxosproto.C_CancelDecision, args.Seqnum, *reply)
    err := json.Unmarshal(r, reply)
    if err != nil {
        fmt.Println(err)
    }
    return nil
}

func (dg *Delegate) DeleteFlight(args * delegateproto.DeleteArgs, reply * delegateproto.DeleteReply) error {
    r := dg.doPush(*args, paxosproto.C_DeleteFlight, args.Seqnum, *reply)
    err := json.Unmarshal(r, reply)
    if err != nil {
        fmt.Println(err)
    }
    return nil
}

func (dg *Delegate) RescheduleFlight(args * delegateproto.RescheduleArgs, reply * delegateproto.RescheduleReply) error {
    r := dg.doPush(*args, paxosproto.C_RescheduleFlight, args.Seqnum, *reply)
    err := json.Unmarshal(r, reply)
    if err != nil {
        fmt.Println(err)
    }
    return nil
}

func (dg *Delegate) QueryFlights(args * delegateproto.QueryArgs, reply * delegateproto.QueryReply) error {
    r := dg.doPush(*args, paxosproto.C_QueryFlights, args.Seqnum, *reply)
    err := json.Unmarshal(r, reply)
    if err != nil {
        fmt.Println(err)
    }
    return nil
}

func (dg *Delegate) AddFlight(args * delegateproto.AddArgs, reply * delegateproto.AddReply) error {
    r := dg.doPush(*args, paxosproto.C_AddFlight, args.Seqnum, *reply)
    err := json.Unmarshal(r, reply)
    if err != nil {
        fmt.Println(err)
    }
    return nil
}

func (dg *Delegate) doPush(args interface{}, t int, seq int, reply interface{}) []byte {
    V := &paxosproto.ValueStruct{}
    buf, _ := json.Marshal(args)
    V.Action = make([]byte, len(buf))
    copy(V.Action, buf)
    V.CoordSeq = seq
    V.Type = t
    rpl := dg.Push(V)
    if rpl.Status == paxosproto.Propose_OK {
        return rpl.Reply
    }
    return nil
}
