package consensus

import (
    "net"
    "net/http"
    "fmt"
    "net/rpc"
    "encoding/json"
    "../config"
    "log"
    "../paxosproto"
    "../airlineserver"
    "sync"
    "strings"
    "strconv"
)

type PaxosEngine struct {
    conf *config.Config
	numNodes int
    log map[int] *paxosproto.ValueStruct
    clients []*rpc.Client
    servers []*paxosproto.NodeStruct
    cur_seq int
    cur_paxos * PaxosInstance
    mutex sync.Mutex
    name string
    in chan * paxosproto.Packet
    out chan * paxosproto.Packet
    brd chan * paxosproto.Packet
	prog chan * paxosproto.ValueStruct 	
	exitCurrentPI chan int
	exitThisEngine chan int
    peerID int
    // airlineserver
    as *airlineserver.AirlineServer

    neth iNetworkHandler
}

//when open a Paxos Engine, the airline name and the peer ID of this server(node) should
// be provided
func NewPaxosEngine(path string, airline_name string, ID int) *PaxosEngine {
    pe := &PaxosEngine{}
    pe.as = airlineserver.NewAirlineServer()
    
//    pe.log = make(map[int] *paxosproto.ValueStruct)

    pe.cur_seq = 0
    peerID := 0
    pe.log = make(map[int] *paxosproto.ValueStruct)
    // read configure file
	conf, _ := config.ReadConfigFile(path)
    pe.conf = conf
    pe.name = airline_name
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
		// client, err := rpc.DialHTTP("tcp",addr)
		// if err != nil {
			// log.Fatal("dialing:", err)
		// }
		// pe.clients[peerID] = client
//		peerID ++
	}
    pe.peerID = ID
	// hostport := pe.servers[ID]		//the host port of this node
	pe.in = make(chan * paxosproto.Packet, 10)
	pe.out = make(chan * paxosproto.Packet, 10)
	pe.brd = make(chan * paxosproto.Packet, 10)
	pe.prog = make(chan * paxosproto.ValueStruct)
	pe.exitCurrentPI = make(chan int, 10)
	pe.exitThisEngine = make(chan int, 10)
	pe.cur_paxos = NewPaxosInstance( ID , pe.cur_seq ,len_list )
	pe.cur_paxos.in = pe.in
	pe.cur_paxos.out = pe.out
	pe.cur_paxos.brd = pe.brd
	pe.cur_paxos.prog = pe.prog
	pe.cur_paxos.log = &pe.log
	pe.cur_paxos.shouldExit = pe.exitCurrentPI
    go pe.cur_paxos.Run()

	go pe.run()
	rpc.Register(pe)
	rpc.HandleHTTP()

    pe.neth.port = airline_servers.UDPPort[ID]
    pe.neth.ReadC = &pe.in
    pe.neth.Listen(pe.neth.port)
    go pe.neth.run()
    fmt.Println("PaxosEngine serving at " + airline_servers.PeersHostPort[ID])

    // fmt.Println(airline_servers.PeersHostPort[ID])
	_, listenport, _ := net.SplitHostPort(airline_servers.PeersHostPort[ID])
	l, e := net.Listen("tcp", fmt.Sprintf(":%s", listenport))
	if e!=nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l,nil)
    return pe
}

func (pe *PaxosEngine) sendMsg(msg *paxosproto.Packet, id int) error {
    // fmt.Println("sending pkt to " + hostport)
    host := strings.Split(pe.conf.Airlines[pe.name].PeersHostPort[id], ":")[0]

    addr, err := net.ResolveUDPAddr("udp4", host + ":" + strconv.FormatInt(int64(pe.conf.Airlines[pe.name].UDPPort[id]), 10))
    if err != nil {
        fmt.Println(err)
        return err
    }
    pe.neth.SendMsg(msg, addr)
    return nil
}

func (pe *PaxosEngine) run() {
    for {
        select{
        case inPkt := <- pe.in:
            pe.cur_paxos.in <- inPkt
        case outPkt := <- pe.out:
            pe.sendMsg(outPkt, outPkt.PeerID)

        case brdMSg := <- pe.brd:
            for i:=0; i<pe.conf.Airlines[pe.name].NumPeers; i++ {
                pe.sendMsg(brdMSg, i)
                // var reply int
                // client, err := rpc.DialHTTP("tcp",addr)/
                // if err == nil {
                    
                // }
            }
        case Va := <- pe.prog:
//            req.reply <- pe.progress(req.Msg.Va)		//reply to where?
//			  pe.progress(req.Msg.Va)
			  fmt.Println("go into progress")
			  pe.progress(Va)
        case <- pe.exitThisEngine:
        	break
        }
    }
}

func (pe *PaxosEngine) progress(V *paxosproto.ValueStruct) interface{} {
    pe.mutex.Lock()
    defer pe.mutex.Unlock()

    reply, err := pe.as.Progress(V)
    // V.reply = make([]byte, len(reply))
    // copy(V.reply, reply)
    fmt.Println("after airline sever's progress")
    pe.log[pe.cur_seq] = V
    pe.cur_seq ++
//    pe.cur_paxos.running = false
	pe.exitCurrentPI <- 1
	pe.in = make(chan * paxosproto.Packet, 10)
	pe.out = make(chan * paxosproto.Packet, 10)
	pe.brd = make(chan * paxosproto.Packet, 10)
	pe.prog = make(chan * paxosproto.ValueStruct, 10)
	pe.exitCurrentPI = make(chan int, 10)
    pe.cur_paxos = NewPaxosInstance(pe.peerID, pe.cur_seq, pe.numNodes)
	pe.cur_paxos.in = pe.in
	pe.cur_paxos.out = pe.out
	pe.cur_paxos.brd = pe.brd
	pe.cur_paxos.prog = pe.prog
	pe.cur_paxos.log = &pe.log
	pe.cur_paxos.shouldExit = pe.exitCurrentPI
    go pe.cur_paxos.Run()
	fmt.Println("after init a new PI")
    if err != nil {
        fmt.Println(err)
    	nop_v := &paxosproto.ValueStruct{}
    	nop_v.Type = paxosproto.C_NOP
    	pe.log[pe.cur_seq] = nop_v
    	return nil
    }
    buf, err := json.Marshal(reply)
    if err != nil {
    	fmt.Println(err)
    }
    V.Reply = make([]byte, len(buf))
    copy(V.Reply, buf)
    
	
	fmt.Println("after successfully write into log")
	
	return reply
}

func (pe * PaxosEngine) checkLog(V * paxosproto.ValueStruct) (bool,int) {
	//no need to lock
	fmt.Println("cehcking log")
	for key, value := range(pe.log) {
		if value.CoordSeq == V.CoordSeq {
			return true, key
		}
	}
	return false,-1
}

func (pe *PaxosEngine) Propose(V * paxosproto.ValueStruct, reply * paxosproto.ReplyStruct) error {
    pe.mutex.Lock()
    defer pe.mutex.Unlock()

    if pe.cur_paxos.Finished {
        reply = nil
        return nil
    }

    fmt.Println("Propose Called")
    fmt.Println(V)

    //check log
    found, index := pe.checkLog(V)
    if found {		//already commited
    	fmt.Println("got log")
    	fmt.Println(index)
    	reply.Status = paxosproto.Propose_OK
    	reply.Type = pe.log[index].Type
    	reply.Reply = make([]byte, len(pe.log[index].Reply))
    	copy(reply.Reply, pe.log[index].Reply)
		return nil
    }

    fmt.Println("No log found")
    
    var Status int
    var Vp * paxosproto.ValueStruct

    for {
        fmt.Println("preparing")
        Status, Vp = pe.cur_paxos.Prepare()
        fmt.Println("Prepare returned")
        if Status == paxosproto.PREPARE_BEHIND {
            //pe.ReqProgress(Vp)
            // pe.progress(Vp)
            pe.prog <- Vp
            reply.Status = paxosproto.Propose_RETRY
            return nil
        } else {
            break
        }
    }
    
    if Status == paxosproto.PREPARE_REJECT {
        fmt.Println("prepare rejected")
        reply.Status = paxosproto.Propose_RETRY
        return nil
    }

    if Vp == nil {
        Vp = V
    }

    fmt.Println("sending accept")
    fmt.Println(Vp)
    OK := pe.cur_paxos.Accept(Vp)
    if OK == -1 {
        fmt.Println("Accept failed")
        reply.Status = paxosproto.Propose_FAIL
        return nil
    }

    pe.cur_paxos.Commit() 
    fmt.Println("{Progessing}")
//    rpl := pe.progress(Vp)
	// pe.prog <- Vp	 
 	reply.Status = paxosproto.Propose_RETRY
 	return nil
    
//    if Vp.CoordSeq != V.CoordSeq {
//    	reply.Status = paxosproto.Propose_RETRY
//    	return nil
//    }
//    
//    fmt.Println(rpl)
//    buf, err := json.Marshal(rpl)
//    if err != nil {
//        fmt.Println(err)
//        return err
//    }
//	fmt.Println("marshal success")
//	fmt.Println(buf)
//	reply.Reply = make([]byte, len(buf))
//    copy(reply.Reply, buf)
//    reply.Type = V.Type
// 	reply.Status = paxosproto.Propose_OK
// 	fmt.Println("returning")
//    return nil
}
