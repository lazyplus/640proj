package consensus

import (
    "../config"
    "../paxosproto"
    "../airlineserver"
    "net"
    "net/http"
    "net/rpc"
    "fmt"
    "encoding/json"
    "log"
    "sync"
    "strings"
    "strconv"
)

type PaxosEngine struct {
    conf *config.Config

    name string
    peerID int

    mutex sync.Mutex
    cur_seq int

    peerAddr [] * net.UDPAddr
    log map[int] *paxosproto.ValueStruct
    in chan * paxosproto.Packet
    out chan * paxosproto.Packet
    brd chan * paxosproto.Packet
	prog chan * paxosproto.ValueStruct
	exitCurrentPI chan int
	exitThisEngine chan int
    
    // components
    as *airlineserver.AirlineServer
    cur_paxos * PaxosInstance
    neth iNetworkHandler
}

//when open a Paxos Engine, the airline name and the peer ID of this server(node) should
// be provided
func NewPaxosEngine(path string, airline_name string, ID int) *PaxosEngine {
    pe := &PaxosEngine{}

    pe.peerID = ID
    pe.name = airline_name

    // read configure file
	conf, _ := config.ReadConfigFile(path)
    pe.conf = conf

	airline_servers, found := conf.Airlines[airline_name]
	if !found {
		return nil
	}

    pe.cur_seq = 0
    pe.log = make(map[int] *paxosproto.ValueStruct)
	pe.in = make(chan * paxosproto.Packet, 10)
    pe.exitThisEngine = make(chan int, 10)

    pe.as = airlineserver.NewAirlineServer()
	pe.progInstance()

    pe.peerAddr = make([] *net.UDPAddr, pe.conf.Airlines[pe.name].NumPeers)
    for i:=0; i<pe.conf.Airlines[pe.name].NumPeers; i++ {
        host := strings.Split(pe.conf.Airlines[pe.name].PeersHostPort[i], ":")[0]
        port := strconv.FormatInt(int64(pe.conf.Airlines[pe.name].UDPPort[i]), 10)
        addr, err := net.ResolveUDPAddr("udp4", host + ":" + port)
        if err != nil {
            fmt.Println(err)
            return nil
        }
        pe.peerAddr[i] = addr
        if i == pe.peerID {
            pe.neth.ReadC = &pe.in
            pe.neth.Listen(host + ":" + port)
            go pe.neth.run()
        }
    }
    
    fmt.Println("PaxosEngine serving at " + airline_servers.PeersHostPort[ID])
	_, listenport, _ := net.SplitHostPort(airline_servers.PeersHostPort[ID])
	l, e := net.Listen("tcp", fmt.Sprintf(":%s", listenport))
	if e!=nil {
		log.Fatal("listen error:", e)
	}
    rpc.Register(pe)
    rpc.HandleHTTP()
	go http.Serve(l,nil)

    go pe.run()

    return pe
}

func (pe *PaxosEngine) progInstance() {
    pe.out = make(chan * paxosproto.Packet, 10)
    pe.brd = make(chan * paxosproto.Packet, 10)
    pe.prog = make(chan * paxosproto.ValueStruct, 10)
    pe.exitCurrentPI = make(chan int, 10)
    pe.cur_paxos = NewPaxosInstance(pe.peerID, pe.cur_seq, pe.conf.Airlines[pe.name].NumPeers)
    pe.cur_paxos.in = make(chan * paxosproto.Packet, 10)
    pe.cur_paxos.out = pe.out
    pe.cur_paxos.brd = pe.brd
    pe.cur_paxos.prog = pe.prog
    pe.cur_paxos.log = &pe.log
    pe.cur_paxos.shouldExit = &pe.exitCurrentPI
    go pe.cur_paxos.Run()
}

func (pe *PaxosEngine) run() {
    for {
        select{
        case inPkt := <- pe.in:
            pe.cur_paxos.in <- inPkt

        case outPkt := <- pe.out:
            pe.neth.SendMsg(outPkt, pe.peerAddr[outPkt.PeerID])

        case brdMSg := <- pe.brd:
            go func () {
                for i:=0; i<pe.conf.Airlines[pe.name].NumPeers; i++ {
                    pe.neth.SendMsg(brdMSg, pe.peerAddr[i])
                }
            } ()
        case Va := <- pe.prog:
			  // fmt.Println("go into progress")
			  pe.progress(Va)

        case <- pe.exitThisEngine:
        	break
        }
    }
}

func (pe *PaxosEngine) progress(V *paxosproto.ValueStruct) {
    pe.mutex.Lock()
    defer pe.mutex.Unlock()

    if V == nil {
        fmt.Println("Shit, nil V")
    }

    found, _ := pe.checkLog(V)
    if !found {
        fmt.Println("progressing", pe.cur_seq)
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
    }

    pe.exitCurrentPI <- 1
    pe.cur_seq ++
    pe.progInstance()
    return
}

func (pe * PaxosEngine) checkLog(V * paxosproto.ValueStruct) (bool,int) {
	//no need to lock
	for key, value := range(pe.log) {
		if value.Type != paxosproto.C_NOP && value.Type == V.Type && value.CoordSeq == V.CoordSeq {
			return true, key
		}
	}
	return false, -1
}

func (pe *PaxosEngine) Propose(V * paxosproto.ValueStruct, reply * paxosproto.ReplyStruct) error {
    pe.mutex.Lock()
    defer pe.mutex.Unlock()

    if pe.cur_paxos.Finished {
        reply.Status = paxosproto.Propose_FAIL
        return nil
    }

    // fmt.Println("Propose Called", pe.cur_seq, V.Type)

    //check log
    found, index := pe.checkLog(V)
    if found {		//already commited
    	// fmt.Println("got log")
    	// fmt.Println(pe.cur_seq, index)
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
        // fmt.Println(pe.cur_seq, "preparing")
        Status, Vp = pe.cur_paxos.Prepare()
        // fmt.Println("Prepare returned")
        if Status == paxosproto.PREPARE_BEHIND {
            pe.cur_paxos.Finished = true
            // fmt.Println("Learning", pe.cur_seq, Vp)
            pe.prog <- Vp
            reply.Status = paxosproto.Propose_LEARN
            return nil
        } else {
            break
        }
    }

    if Status == paxosproto.PREPARE_REJECT {
        fmt.Println(pe.cur_seq, "prepare rejected")
        reply.Status = paxosproto.Propose_FAIL
        return nil
    }

    if Vp == nil {
        Vp = V
    }

    // fmt.Println(pe.cur_seq, "sending accept")

    OK := pe.cur_paxos.Accept(Vp)
    if OK != paxosproto.ACCEPT_OK {
        fmt.Println(pe.cur_seq, "Accept failed")
        reply.Status = paxosproto.Propose_FAIL
        return nil
    }

    pe.cur_paxos.Commit(Vp)
    // fmt.Println(pe.cur_seq, "committing")
    pe.cur_paxos.Finished = true
    pe.prog <- Vp
 	reply.Status = paxosproto.Propose_RETRY
 	return nil
}
