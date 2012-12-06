package delegateimpl

import (
	"time"
    "fmt"
    "net/rpc"
    "../paxosproto"
    "../config"
    "../delegateproto"
    "encoding/json"
    "math/rand"
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
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for index := r.Intn(dg.conf.NumPeers); ; index = (index + 1) % dg.conf.NumPeers {
		reply := &paxosproto.ReplyStruct{}
        client, err := rpc.DialHTTP("tcp", dg.conf.PeersHostPort[index])
        if err != nil {
            continue
        }
        retc := make(chan error);
        go func (cli *rpc.Client, V *paxosproto.ValueStruct, reply *paxosproto.ReplyStruct)  {
            err := client.Call("PaxosEngine.Propose", V, reply)
            retc <- err
        } (client, V, reply)
        timer := time.NewTimer(5 * time.Second)
        select{
        case <- timer.C:
            // time out
            fmt.Println("call Propose time out")
            client.Close()
            timer.Stop()
        case err := <- retc:
            timer.Stop()
            if err != nil {
                fmt.Println("call Propose err: ", err)
                client.Close()
                continue
            } else {
                client.Close()
                if reply.Status == paxosproto.Propose_OK {
                    return reply
                }
                if reply.Status == paxosproto.Propose_LEARN {
                    index = index + dg.conf.NumPeers - 1
                    // continue
                }
            }
        }
        time.Sleep(5 * time.Millisecond)
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
