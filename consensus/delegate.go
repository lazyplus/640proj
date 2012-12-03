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
	for e := airline_server_list.Front(); e != nil; e = e.Next() {
		addr := e.Value
		dg.servers[i] = addr
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			cli[i] = client
			i++
		}
		//else : give up this node
	}
	dg.Port = port
	return dg
}

func (dg *Delegate) Push(V ValueStruct) interface{} {
    index := 0

    for ; ; index = (index + 1) % len(dg.servers) {
        var reply interface{}
        dg.cli[index].Call("AirlineServerRPC.Propose", V, reply)
        if reply != nil {
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
    reply = dg.Push(V)
    return nil
}

func (dg *Delegate) QueryFlights(args * QueryArgs, reply * QueryReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_QueryFlights
    reply = dg.Push(V)
    return nil
}

func (dg *Delegate) PrepareBookFlight(args * BookArgs, reply * BookReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_PrepareBookFlight
    reply = dg.Push(V)
    return nil
}

func (dg *Delegate) BookDecision(args * DecisionArgs, reply * DecisionReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_BookDecision
    reply = dg.Push(V)
    return nil
}

func (dg *Delegate) CancelDecision(args * DecisionArgs, reply * DecisionReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_CancelDecision
    reply = dg.Push(V)
    return nil
}

func (dg *Delegate) DeleteFlight(args * DeleteArgs, reply * DeleteReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_DeleteFlight
    reply = dg.Push(V)
    return nil
}

func (dg *Delegate) RescheduleFlight(args * RescheduleArgs, reply * RescheduleReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_RescheduleFlight
    reply = dg.Push(V)
    return nil
}

func (dg *Delegate) AddFlight(args * AddArgs, reply * AddReply) error {
    V := ValueStruct{}
    V.action = args
    V.port = dg.Port
    V.CoordSeq = args.seqnum
    V.Type = c_AddFlight
    reply = dg.Push(V)
    return nil
}
