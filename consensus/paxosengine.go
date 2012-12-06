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
    // "time"
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
    behind bool
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
	pe.prog = make(chan * paxosproto.ValueStruct, 10)
	pe.exitCurrentPI = make(chan int, 10)
	pe.exitThisEngine = make(chan int, 10)
	pe.cur_paxos = NewPaxosInstance( ID , pe.cur_seq ,len_list )
	pe.cur_paxos.in = pe.in
	pe.cur_paxos.out = pe.out
	pe.cur_paxos.brd = pe.brd
	pe.cur_paxos.prog = pe.prog
	pe.cur_paxos.log = &pe.log
    pe.cur_paxos.mutex = &pe.mutex
	pe.cur_paxos.shouldExit = &pe.exitCurrentPI
    go pe.cur_paxos.Run()

    pe.behind = false
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

    addr, err := net.ResolveUDPAddr("udp", host + ":" + strconv.FormatInt(int64(pe.conf.Airlines[pe.name].UDPPort[id]), 10))
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
            }

        case Va := <- pe.prog:
//            req.reply <- pe.progress(req.Msg.Va)		//reply to where?
//			  pe.progress(req.Msg.Va)
			  // fmt.Println("go into progress")
			  pe.progress(Va)
        case <- pe.exitThisEngine:
        	break
        }
    }
}

func (pe *PaxosEngine) progress(V *paxosproto.ValueStruct) interface{} {
    // fmt.Println("Getting Lock")
    pe.mutex.Lock()
    defer pe.mutex.Unlock()
    // fmt.Println("Got Lock")
    if V == nil {
        fmt.Println("Warning: Progressing nil")
    }

    fmt.Println("progressing", pe.cur_seq, V)
    reply, err := pe.as.Progress(V)

    if err == nil {
        buf, err := json.Marshal(reply)
        if err != nil {
            fmt.Println(err)
        }
        V.Reply = make([]byte, len(buf))
        copy(V.Reply, buf)
        pe.log[pe.cur_seq] = V
    } else {
        nop_v := &paxosproto.ValueStruct{}
        nop_v.Type = paxosproto.C_NOP
        pe.log[pe.cur_seq] = nop_v
    }
    
    pe.cur_seq ++
	pe.exitCurrentPI <- 1
	pe.out = make(chan * paxosproto.Packet, 10)
	pe.brd = make(chan * paxosproto.Packet, 10)
	pe.prog = make(chan * paxosproto.ValueStruct, 10)
	pe.exitCurrentPI = make(chan int, 10)
    pe.cur_paxos = NewPaxosInstance(pe.peerID, pe.cur_seq, pe.numNodes)
	pe.cur_paxos.in = make(chan * paxosproto.Packet, 10)
	pe.cur_paxos.out = pe.out
	pe.cur_paxos.brd = pe.brd
	pe.cur_paxos.prog = pe.prog
	pe.cur_paxos.log = &pe.log
    pe.cur_paxos.mutex = &pe.mutex
	pe.cur_paxos.shouldExit = &pe.exitCurrentPI
    go pe.cur_paxos.Run()
	// fmt.Println("after init a new PI")
    

	// fmt.Println("after successfully write into log")
	
	return reply
}

func (pe * PaxosEngine) checkLog(V * paxosproto.ValueStruct) (bool,int) {
	//no need to lock
	// fmt.Println("cehcking log")
	for key, value := range(pe.log) {
		if value.Type == V.Type && value.CoordSeq == V.CoordSeq {
			return true, key
		}
	}
	return false,-1
}

func (pe *PaxosEngine) Propose(V * paxosproto.ValueStruct, reply * paxosproto.ReplyStruct) error {
    // ProgV, ok := <- pe.prog
    // if ok {
    //     pe.prog <- ProgV
    //     time.Sleep(10 * time.Millisecond)
    // }

    pe.mutex.Lock()
    defer pe.mutex.Unlock()

    if pe.cur_paxos.Finished {
        reply.Status = paxosproto.Propose_FAIL
        return nil
    }

    fmt.Println("Propose Called", pe.cur_seq, V.Type)

    //check log
    found, index := pe.checkLog(V)
    if found {		//already commited
    	// fmt.Println("got log")
    	fmt.Println(pe.cur_seq, index)
    	reply.Status = paxosproto.Propose_OK
    	reply.Type = pe.log[index].Type
    	reply.Reply = make([]byte, len(pe.log[index].Reply))
    	copy(reply.Reply, pe.log[index].Reply)
		return nil
    }

    // fmt.Println("No log found")
    
    var Status int
    var Vp * paxosproto.ValueStruct

    for {
        fmt.Println(pe.cur_seq, "preparing")
        Status, Vp = pe.cur_paxos.Prepare()
        // fmt.Println("Prepare returned")
        if Status == paxosproto.PREPARE_BEHIND {
            pe.cur_paxos.Finished = true
            pe.behind = true
            fmt.Println("Learning", pe.cur_seq, Vp)
            pe.prog <- Vp
            reply.Status = paxosproto.Propose_LEARN
            return nil
        } else {
            break
        }
    }
    
    pe.behind = false
    if Status == paxosproto.PREPARE_REJECT {
        fmt.Println(pe.cur_seq, "prepare rejected")
        reply.Status = paxosproto.Propose_FAIL
        return nil
    }

    if Vp == nil {
        Vp = V
    }

    fmt.Println(pe.cur_seq, "sending accept")

    OK := pe.cur_paxos.Accept(Vp)
    if OK == -1 {
        fmt.Println(pe.cur_seq, "Accept failed")
        reply.Status = paxosproto.Propose_FAIL
        return nil
    }
    if OK == paxosproto.ACCEPT_REJECT {
        reply.Status = paxosproto.Propose_FAIL
        return nil
    }

    pe.cur_paxos.Commit(Vp)
    fmt.Println(pe.cur_seq, "committing")
    pe.cur_paxos.Finished = true
    pe.prog <- Vp
 	reply.Status = paxosproto.Propose_RETRY
 	return nil
}
