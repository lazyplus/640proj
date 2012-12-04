package consensus

import (
    "net/http"
    "net/rpc"
)

func NewDelegate(path string, airline_name string, port string) *Delegate {
    dg := &Delegate{}
	i := 0
	cg, _ := ReadConfigFile(path)
	airline_server_list, found := cg.AirlineAddr[airline_name]
	if !found {
		return nil
	}
	len_list := len(airline_server_list)
	dg.numServers = len_list
	dg.servers = make([]string,len_list)
	dg.cli = make([]*rpc.Client,len_list)
	for e := airline_server_list.Front(); e != nil; e = e.Next(), i++ {
		addr := e.Value
		dg.servers[i] = addr
		dg.cli[i] = nil
	}
	dg.Port = port
	return dg
}

func (dg *Delegate) Push(V ValueStruct) interface{} {
    index := 0

    for ; ; index = (index + 1) % len(dg.servers) {
        var reply interface{}
        client, err := rpc.DialHTTP("tcp", dg.servers)
        if err != nil {
            continue
        }
        // TODO: prepare RPC client
        err = client.Call("AirlineServerRPC.Propose", V, reply)
        if err != nil {
            client.Close()
            continue
        }
        if reply != nil {
            client.Close()
            return reply
        }
        timer.Sleep(time.Second)
    }
}

func (dg *Delegate) PrepareCancelFlight(args * BookArgs, reply * BookReply) error {
    // make up Value
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_PrepareCancelFlight
    *reply = dg.Push(V)
    return nil
}

func (dg *Delegate) QueryFlights(args * QueryArgs, reply * QueryReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_QueryFlights
    *reply = dg.Push(V)
    return nil
}

func (dg *Delegate) PrepareBookFlight(args * BookArgs, reply * BookReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_PrepareBookFlight
    *reply = dg.Push(V)
    return nil
}

func (dg *Delegate) BookDecision(args * DecisionArgs, reply * DecisionReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_BookDecision
    *reply = dg.Push(V)
    return nil
}

func (dg *Delegate) CancelDecision(args * DecisionArgs, reply * DecisionReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_CancelDecision
    *reply = dg.Push(V)
    return nil
}

func (dg *Delegate) DeleteFlight(args * DeleteArgs, reply * DeleteReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_DeleteFlight
    *reply = dg.Push(V)
    return nil
}

func (dg *Delegate) RescheduleFlight(args * RescheduleArgs, reply * RescheduleReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_RescheduleFlight
    *reply = dg.Push(V)
    return nil
}

func (dg *Delegate) AddFlight(args * AddArgs, reply * AddReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_AddFlight
    *reply = dg.Push(V)
    return nil
}
