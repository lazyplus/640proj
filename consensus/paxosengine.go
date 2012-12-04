package consensus

import (
    "net"
    "net/http"
    "net/rpc"
    "../config"
    "log"
    "../paxosproto"
    "../airlineserver"
    "sync"
)


type PaxosEngine struct {
	numNodes int
    log map[int] *paxosproto.ValueStruct
    clients []*rpc.Client
    servers []*paxosproto.NodeStruct
    cur_seq int
    cur_paxos * PaxosInstance
    mutex sync.Mutex
    in chan * paxosproto.Packet
    out chan * paxosproto.Packet
    brd chan * paxosproto.Packet
	prog chan * paxosproto.Packet
	exitCurrentPI chan int
	exitThisEngine chan int
    peerID int
    // airlineserver
    as *airlineserver.AirlineServer
}

//when open a Paxos Engine, the airline name and the peer ID of this server(node) should
// be provided
func NewPaxosEngine(path string, airline_name string, ID int) *PaxosEngine {
    pe := &PaxosEngine{}
    pe.as = airlineserver.NewAirlineServer()

    pe.cur_seq = 0
    peerID := 0
    pe.log = make(map[int] *paxosproto.ValueStruct)
    // read configure file
	conf, _ := config.ReadConfigFile(path)
	airline_servers, found := conf.Airlines[airline_name]
	if !found {
		return nil
	}	
	airline_server_list := airline_servers.PeersHostPort
	len_list := len(airline_server_list)
	pe.numNodes = len_list
	pe.servers = make([]*paxosproto.NodeStruct,len_list)
	pe.clients = make([]*rpc.Client,len_list)
	//for e := airline_server_list.Front(); e != nil; e = e.Next() {
	for peerID=0; peerID < len_list; peerID ++ {
//		addr := e.Value
		addr := airline_server_list[peerID]
		pe.servers[peerID] = &paxosproto.NodeStruct{addr,peerID}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		pe.clients[peerID] = client
//		peerID ++
	}
	hostport := pe.servers[ID]		//the host port of this node
	pe.in = make(chan * paxosproto.Packet)
	pe.out = make(chan * paxosproto.Packet)
	pe.brd = make(chan * paxosproto.Packet)
	pe.prog = make(chan * paxosproto.Packet)
	pe.exitCurrentPI = make(chan int)
	pe.exitThisEngine = make(chan int)
	pe.cur_paxos = NewPaxosInstance( ID , pe.cur_seq ,len_list )
	pe.cur_paxos.in = pe.in
	pe.cur_paxos.out = pe.out
	pe.cur_paxos.brd = pe.brd
	pe.cur_paxos.prog = pe.prog
	pe.cur_paxos.log = pe.log
	pe.cur_paxos.shouldExit = pe.exitCurrentPI
	go pe.Run()
	rpc.Register(pe)
	rpc.HandleHTTP()
	_, listenport, _ := net.SplitHostPort(hostport.Port)
	l, e := net.Listen("tcp",listenport)
	if e!=nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l,nil)
    return pe
}

//use RPC to coordinate between engines
//pass the msg to the channel
func (pe * PaxosEngine) receiveRPC ( args * paxosproto.Packet, reply * int) {
	pe.in <- args
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
					// var reply int                   
					pe.clients[i].Go("RPCReceiver.receiveRPC",outPkt, nil,nil)	                   
                    break
                }
            }
        case brdMSg := <- pe.brd:
            for i:=0; i<len(pe.servers); i++ {
                // var reply int
                pe.clients[i].Go("RPCReceiver.receiveRPC",brdMSg, nil, nil)
            }
        case req := <- pe.prog:
//            req.reply <- pe.progress(req.Msg.Va)		//reply to where?
			  pe.progress(req.Msg.Va)
        case <- pe.exitThisEngine:
        	break
        }
    }
}

func (pe *PaxosEngine) progress(V paxosproto.ValueStruct) interface{} {
    reply, err := pe.as.Progress(V)
    pe.cur_seq ++
    pe.exitCurrentPI <- 1
    pe.cur_paxos = NewPaxosInstance(pe.peerID, pe.cur_seq, pe.numNodes)

    if err != nil {
    	nop_v := &paxosproto.ValueStruct{}
    	nop_v.Type = paxosproto.C_NOP
    	pe.log[pe.cur_seq] = nop_v
    	return nil
    }
	pe.log[pe.cur_seq] = &V
	return reply
}

func (pe * PaxosEngine) CheckLog(V * paxosproto.ValueStruct) (bool,int) {
	//no need to lock
	for i:=0;i<len(pe.log);i++ {
		if pe.log[i].CoordSeq == V.CoordSeq {
			return true,i
		}
	}
	return false,-1
}

func (pe *PaxosEngine) Propose(V * paxosproto.ValueStruct, reply * paxosproto.ReplyStruct) {
    pe.mutex.Lock()
    defer pe.mutex.Unlock()

    //check log
    found, index := pe.CheckLog(V)
    if found {		//already commited
    	reply.Status = paxosproto.Propose_OK
    	reply.Type = pe.log[index].Type
    	reply.Reply = pe.log[index].Reply
		return
    }
    
    var Status int
    var Vp * paxosproto.ValueStruct

    for {
        Status, Vp = pe.cur_paxos.Prepare()
        if Status == paxosproto.PREPARE_BEHIND {
            //pe.ReqProgress(Vp)
            pe.progress(*Vp)
        } else {
            break
        }
    }
    
    if Status == paxosproto.PREPARE_REJECT {
        reply.Status = paxosproto.Propose_RETRY
        return 
    }

    if Vp == nil {
        Vp = V
    }

    OK := pe.cur_paxos.Accept()
    if OK == -1 {
        reply.Status = paxosproto.Propose_FAIL
        return 
    }

    pe.cur_paxos.Commit() 
	reply.Reply = pe.progress(*Vp)
 	reply.Status = paxosproto.Propose_OK
    return 
}
