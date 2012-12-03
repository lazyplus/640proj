package consensus

import (
    "../paxosinstance"
    "net"
    "net/http"
    "net/rpc"
    "../config"
    "log"
)

//when open a Paxos Engine, the airline name and the peer ID of this server(node) should
// be provided
func NewPaxosEngine(path String, airline_name string, ID int) *PaxosEngine {
    pe := &PaxosEngine{}
    pe.cur_seq = 0
    pe.log = make(map[int] *ValueStruct)
    // read configure file
	conf, _ := config.ReadConfigFile(path)
	airline_server_list, found := conf.AirlineAddr[airline_name]
	if !found {
		return nil
	}	
	len_list = len(airline_server_list)
	pe.num = len_list
	pe.servers = make([]*NodeStruct,len_list)
	pe.clients = make([]*rpc.Client)
	peerID := 0
	for e := airline_server_list.Front(); e != nil; e = e.Next() {
		addr := e.Value
		pe.servers[peerID] = &NodeStruct{addr,peerID}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		pe.clients[peerID] = client
		peerID ++
	}
	hostport := pe.servers[ID]
    // init netHandler
	pe.in = make(chan * Packet)
	pe.out = make(chan * Packet)
	pe.brd = make(chan * MsgStruct)
	pe.prog = make(chan * Packet)
	pe.exitCurrentPI = make(chan int)
	pe.cur_paxos = NewPaxosInstance( ID , cur_seq ,len_list )
	pe.cur_paxos.in = pe.in
	pe.cur_paxos.out = pe.out
	pe.cur_paxos.brd = pe.brd
	pe.cur_paxos.prog = pe.prog
	pe.cur_paxos.log = pe.log
	pe.cur_paxos.shouldExit = pe.exitCurrentPI
	pe.RPCReceiver = new(RPCStruct)
	pe.RPCReceiver.in = pe.in
	rpc.Register(pe.RPCReceiver)
	rpc.HandleHTTP()
	_, listenport, _ := net.SplitHostPort(hostport)
	l, e := net.Listen("tcp",listenport)
	if e!=nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l,nil)
    return pe
}

//use RPC to coordinate between engines
//pass the msg to the channel
func (rpcrecv * RPCReceiver) receiveRPC ( args * Packet, reply * int) {
	rpcrecv.in <- args.Msg
	reply = nil
}

func (pe *PaxosEngine) Run() {
    for {
        select{
        case inPkt := <- pe.in:
            pe.cur_paxos.in <- inPkt
        case outPkt := <- pe.out:
            for i:=0; i<len(pe.servers); i++ {
                if pe.servers[i].ID == outPkt.PeerID {
                    //pe.networkHandler.sendMsg(outMsg, servers.Addr)
					var reply int                   
					pe.clients[i].Go("RPCReceiver.receiveRPC",outPkt, nil,nil)
					                   
                    break
                }
            }
        }
        case brdMSg := <- pe.brd:
            for i:=0; i<len(pe.servers); i++ {
                //pe.networkHandler.sendMsg(brdMsg, servers.Addr)
                var reply int
                pe.clients[i].Go("RPCReceiver.receiveRPC",brdMsg, nil, nil)
            }
        case req := <- pe.prog:
            req.reply <- pe.progress(req.Msg.Va)
        }
    }
}

func (pe *PaxosEngine) progress(V ValueStruct) {
    pe.as.Progress(V)
    pe.log[cur_seq] = V
    pe.cur_seq ++
    pe.exitCurrentPI <- 1
    pe.cur_paxos = NewPaxosInstance(cur_seq)
}



func (pe *PaxosEngine) Propose() {
    pe.mutex.Lock()
    defer pe.mutex.Unlock()

    //check log
    reply.Result = pe.as.CheckLog(V)
    if reply.Result != nil {
        reply.Status = OK
        return nil
    }

    var Status int
    var Vp ValueStruct

    for {
        Status, Vp = pe.cur_paxos.prepare()
        if Status == PREPARE_BEHIND {
            pe.ReqProgress(Vp)
        } else {
            break
        }
    }
    
    if Status == PREPARE_REJECT {
        reply.Status = RETRY
        return nil
    }

    if Vp == nil {
        Vp = args.V
    }

    OK = pe.cur_paxos.accept(Vp)
    if OK == -1 {
        reply.Status = FAILED
        return nil
    }

    pe.cur_paxos.commit(Vp) 

    reply.Result = pe.ReqProgress(Vp)
    reply.Status = FAILED
    return nil
}
