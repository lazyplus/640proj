package consensus

import (
	"../paxosproto"
	"time"
)

func generate_random_number(PeerID int, numNodes int, Nh int) int {
    return (Nh / 10 + 1) * 10 + PeerID
}

type PaxosInstance struct {
	log *map[int] *paxosproto.ValueStruct
    out chan * paxosproto.Packet
    brd chan * paxosproto.Packet
    in chan * paxosproto.Packet	
    prog chan * paxosproto.ValueStruct
    prepareCh chan * paxosproto.MsgStruct
    acceptCh chan * paxosproto.MsgStruct
    shouldExit *chan int

	Na int
	Nh int
	Myn int
	Va *paxosproto.ValueStruct
    seq int
    PeerID int
    numNodes int

    PreaccepteNodes int
    PrefailNodes int
    AcpacceptedNodes int
    AcpfailNodes int

    running bool
    Finished bool
}

func NewPaxosInstance(peerID int, seq int, numNodes int) *PaxosInstance {
    pi := &PaxosInstance{}
    pi.seq = seq
    pi.PeerID = peerID
    pi.numNodes = numNodes
    pi.prepareCh = make(chan * paxosproto.MsgStruct, 10)
    pi.acceptCh = make(chan * paxosproto.MsgStruct, 10)
    pi.running = true
    return pi
}

func (pi *PaxosInstance) Run() {
    for pi.running {
        select{
        case <- (*pi.shouldExit):
        	pi.running = false
        	break

        case inPkt := <- pi.in:
			if pi.running == false {
				break
			}

            if inPkt.Msg.Type == paxosproto.PREPARE {
                pi.handlePrepare(inPkt)
                break
            }
            
            if inPkt.Msg.Type == paxosproto.ACCEPT {
                pi.handleAccept(inPkt)
                break
            }

            if inPkt.Msg.Seq != pi.seq {
                break
            }

            switch(inPkt.Msg.Type){
            case paxosproto.PREPARE_OK:
				pi.PreaccepteNodes ++
                // fmt.Println("Get prepare OK")
                if pi.Na < inPkt.Msg.Na {
                    pi.Na = inPkt.Msg.Na
                    pi.Va = inPkt.Msg.Va
                }
				if (pi.PreaccepteNodes > (pi.numNodes/2)) {
                    // msg := &paxosproto.MsgStruct{}
                    // msg.Type = paxosproto.PREPARE_OK
                    // msg.Na = pi.Na
                    // msg.Va = pi.Va
                    // pi.prepareCh <- msg
                    pi.prepareCh <- inPkt.Msg
				}
            case paxosproto.PREPARE_REJECT:
                // fmt.Println("Get prepare rejected")
                pi.PrefailNodes ++
                if (pi.PrefailNodes > (pi.numNodes/2))  {
                     pi.prepareCh <- inPkt.Msg
                }
            case paxosproto.PREPARE_BEHIND:
                pi.Va = inPkt.Msg.Va
                pi.prepareCh <- inPkt.Msg
            
            case paxosproto.ACCEPT_OK:
                // fmt.Println("Get ACCEPT_OK")
                pi.AcpacceptedNodes ++
                if (pi.AcpacceptedNodes > (pi.numNodes/2)) {
                	pi.acceptCh <- inPkt.Msg
                }
            case paxosproto.ACCEPT_REJECT:
                // fmt.Println("Get ACCEPT_REJECT")
				pi.AcpfailNodes ++
				if (pi.AcpfailNodes > (pi.numNodes/2)){
                	pi.acceptCh <- inPkt.Msg
                }
            case paxosproto.COMMIT:
                pi.handleCommit(inPkt)
            }
        }
    }
}

func (pi *PaxosInstance) initPkt () *paxosproto.Packet {
	p := &paxosproto.Packet{}
	p.Msg = &paxosproto.MsgStruct{}
	p.PeerID = pi.PeerID
	p.Msg.Seq = pi.seq
	return p
}

func (pi *PaxosInstance) handlePrepare(pkt * paxosproto.Packet) {
	msg := pkt.Msg
    // fmt.Println("Get prepare, my seq")
    // fmt.Println(pi.seq)
    // fmt.Println("msg seq is ")
    // fmt.Println(pkt.Msg.Seq)
    if pkt.Msg.Seq < pi.seq {
        // reply prepare behind
        newPkt := pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Seq = pkt.Msg.Seq
        newPkt.Msg.Type = paxosproto.PREPARE_BEHIND
        newPkt.Msg.Va = (*pi.log)[pkt.Msg.Seq]
        pi.out <- newPkt
        return
    } else if pkt.Msg.Seq > pi.seq {
        return
    }

    // fmt.Println("get prepare", msg.Na, pi.Nh)
    if msg.Na < pi.Nh {
        // reply prepare rejected
 		newPkt := pi.initPkt()
 		newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = paxosproto.PREPARE_REJECT
        pi.out <- newPkt
    }else {
        pi.Nh = msg.Na
        // reply prepare OK, Va
        newPkt := pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = paxosproto.PREPARE_OK
        newPkt.Msg.Va = pi.Va
        newPkt.Msg.Na = pi.Na
        pi.out <- newPkt
    }
}

func (pi *PaxosInstance) handleAccept(pkt * paxosproto.Packet) {
    // fmt.Println("Handling Accept")
    // fmt.Println(pkt)
	msg := pkt.Msg

    if pkt.Msg.Seq < pi.seq {
        // reply prepare behind
        newPkt := pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Seq = pkt.Msg.Seq
        newPkt.Msg.Type = paxosproto.ACCEPT_REJECT
        pi.out <- newPkt
        return
    } else if pkt.Msg.Seq > pi.seq {
        return
    }

    if msg.Na < pi.Nh {
        // reply accept rejected
        newPkt := pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = paxosproto.ACCEPT_REJECT
        // fmt.Println("rejecting")
        pi.out <- newPkt
    }else{
        pi.Na = msg.Na
        pi.Va = msg.Va
        pi.Nh = msg.Na
        // reply accept OK
        newPkt := pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = paxosproto.ACCEPT_OK
        // fmt.Println("accepting")
        pi.out <- newPkt
    }
}

func (pi *PaxosInstance) handleCommit(pkt * paxosproto.Packet) {
	//receive commit, notify paxosengine to record the log and take action
    if pi.Finished {
        return
    }
    // fmt.Println(pi.seq, "be committed")
    pi.Finished = true
	pi.prog <- pkt.Msg.Va
}

func (pi *PaxosInstance) Prepare() (int, * paxosproto.ValueStruct) {
    pi.PreaccepteNodes = 0
    pi.PrefailNodes = 0
    pi.Myn = generate_random_number(pi.PeerID,pi.numNodes, pi.Nh)

    msg := &paxosproto.MsgStruct{}
    msg.Type = paxosproto.PREPARE
    msg.Seq = pi.seq
    msg.Na = pi.Myn
    newPkt := &paxosproto.Packet{pi.PeerID,msg}
    pi.brd <- newPkt

    timer := time.NewTimer(200 * time.Millisecond)
    p := &paxosproto.MsgStruct{}
    select{
    case <- timer.C:
        p.Type = paxosproto.PREPARE_REJECT
    case p = <- pi.prepareCh:
        timer.Stop()
    }

    var state int
    var v * paxosproto.ValueStruct = nil
    switch p.Type {
    	case paxosproto.PREPARE_OK:
    		state = paxosproto.PREPARE_OK
    		v = pi.Va
    	case paxosproto.PREPARE_REJECT:
 			state = paxosproto.PREPARE_REJECT   	
    	case paxosproto.PREPARE_BEHIND:
    		state = paxosproto.PREPARE_BEHIND
    		v = pi.Va
    }
    return state, v
}

func (pi * PaxosInstance) Accept(v *paxosproto.ValueStruct) int {
    pi.AcpacceptedNodes = 0
    pi.AcpfailNodes = 0

	msg := &paxosproto.MsgStruct{}
	msg.Type = paxosproto.ACCEPT
	msg.Na = pi.Myn
	msg.Va = v
    msg.Seq = pi.seq
	newPkt := &paxosproto.Packet{pi.PeerID,msg}
	pi.brd <- newPkt

    timer := time.NewTimer(200 * time.Millisecond)
    p := &paxosproto.MsgStruct{}
    select{
    case <- timer.C:
        p.Type = paxosproto.ACCEPT_REJECT
    case p = <- pi.acceptCh:
        timer.Stop()
    }

	switch p.Type {
		case paxosproto.ACCEPT_OK:
			return paxosproto.ACCEPT_OK	
		case paxosproto.ACCEPT_REJECT:
			return paxosproto.ACCEPT_REJECT
	}
	return paxosproto.ACCEPT_REJECT
}

func (pi * PaxosInstance) Commit(v *paxosproto.ValueStruct) {
	msg := &paxosproto.MsgStruct{}
	msg.Type = paxosproto.COMMIT
	msg.Va = v
    msg.Seq = pi.seq
	newPkt := &paxosproto.Packet{pi.PeerID,msg}
	pi.brd <- newPkt
}
