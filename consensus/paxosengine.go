package consensus

import (
    "../paxosinstance"
    "net"
)

type PaxosEngine struct {
    log map[int] *logStruct

    cur_seq int
    cur_paxos PaxosInstance
    mutex sync.Mutex

    // network
    netHandler iNetworkHandler
}

func NewPaxosEngine(path String) *PaxosEngine {
    pe := &PaxosEngine{}

    pe.cur_seq = 0
    pe.cur_paxos = NewPaxosInstance(0)

    // read configure file

    // init netHandler

    return pe
}

type ValueStruct struct {
    CoordSeq int
    Type int
    Action interface{}
    Host string
}

type Packet struct {
    PeerID int
    Msg MsgStruct
}

func (pe *PaxosEngine) Run() {
    for {
        select{
        case inPkt := <- pe.in:
            pe.cur_paxos.in <- inPkt
        case outPkt := <- pe.out:
            for i:=0; i<len(pe.servers); i++ {
                if pe.servers[i].ID == outPkt.PeerID {
                    pe.networkHandler.sendMsg(outMsg, servers.Addr)
                    break
                }
            }
        case brdMSg := <- pe.brd:
            for i:=0; i<len(pe.servers); i++ {
                pe.networkHandler.sendMsg(brdMsg, servers.Addr)
            }
        case req := <- pe.prog:
            req.reply <- pe.progress(req.V)
        }
    }
}

func (pe *PaxosEngine) progress(V ValueStruct) {
    pe.as.Progress(V)
    pe.cur_seq ++
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
    if !OK {
        reply.Status = FAILED
        return nil
    }

    OK = pe.cur_paxos.commit(Vp)

    reply.Result = pe.ReqProgress(Vp)
    reply.Status = FAILED
    return nil
}
