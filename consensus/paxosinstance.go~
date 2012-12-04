package consensus

import (
	"math/rand"
	"../paxosproto"
	"time"
)

func generate_random_number(PeerID int, numNodes int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return (r.Intn(10)*numNodes + PeerID)
}

type PaxosInstance struct {
	log map[int] *paxosproto.ValueStruct
    out chan * paxosproto.Packet
    brd chan * paxosproto.Packet
    in chan * paxosproto.Packet	
    prog chan * paxosproto.Packet
    prepareCh chan * paxosproto.MsgStruct
    acceptCh chan * paxosproto.MsgStruct
    shouldExit chan int
	Na int
	Nh int
	Myn int
	Va paxosproto.ValueStruct
    seq int
    PeerID int
    numNodes int
    PreaccepteNodes int
    PrefailNodes int
    AcpacceptedNodes int
    AcpfailNodes int
    running bool
}

func NewPaxosInstance(peerID int, seq int, numNodes int) *PaxosInstance {
    pi := &PaxosInstance{}
    pi.seq = seq
    pi.PeerID = peerID
    pi.numNodes = numNodes
    pi.Myn = generate_random_number(peerID,numNodes)
    pi.prepareCh = make(chan * paxosproto.MsgStruct)
    pi.acceptCh = make(chan * paxosproto.MsgStruct)
    pi.running = true
    go pi.Run()
    return pi
}

func (pi *PaxosInstance) Run() {
	// finishPrepare := false
	// finishAccept := false
	// finishCommit := false
	
    for pi.running {
        select{
        case inPkt := <- pi.in:
            if inPkt.Msg.Type == paxosproto.PREPARE {
                pi.handlePrepare(*inPkt)
                break
            }

            if inPkt.Msg.Seq != pi.seq {
                break
            }

            switch(inPkt.Msg.Type){
            case paxosproto.PREPARE_OK:
				pi.PreaccepteNodes ++
                if pi.Na < inPkt.Msg.Na {
                    pi.Na = inPkt.Msg.Na
                    pi.Va = inPkt.Msg.Va
                }
				if (pi.PreaccepteNodes > (pi.numNodes/2)) {
					pi.prepareCh <- inPkt.Msg
				}
            case paxosproto.PREPARE_REJECT:
                pi.PrefailNodes ++
                if (pi.PrefailNodes > (pi.numNodes/2))  {
 					pi.prepareCh <- inPkt.Msg
                }
            case paxosproto.PREPARE_BEHIND:
                pi.Va = inPkt.Msg.Va
				pi.prepareCh <- inPkt.Msg
            case paxosproto.ACCEPT:
                pi.handleAccept(*inPkt)
            case paxosproto.ACCEPT_OK:
                pi.AcpacceptedNodes ++
                if (pi.AcpacceptedNodes > (pi.numNodes/2)) {
                	pi.acceptCh <- inPkt.Msg
                }
            case paxosproto.ACCEPT_REJECT:
				pi.AcpfailNodes ++
				if (pi.AcpacceptedNodes > (pi.numNodes/2)){
                	pi.acceptCh <- inPkt.Msg
                }
            case paxosproto.COMMIT:
                pi.handleCommit()
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

func (pi *PaxosInstance) handlePrepare(pkt paxosproto.Packet) {
	msg := pkt.Msg
    if pkt.Msg.Seq < pi.seq {
        // reply prepare behind
        newPkt := pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = paxosproto.PREPARE_BEHIND
        newPkt.Msg.Va = *pi.log[pkt.Msg.Seq]
        pi.out <- newPkt
    }
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

func (pi *PaxosInstance) handleAccept(pkt paxosproto.Packet) {
	msg := pkt.Msg
    if msg.Na < pi.Nh {
        // reply accept rejected
        newPkt := pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = paxosproto.ACCEPT_REJECT
        pi.out <- newPkt
    }else{
        pi.Na = msg.Na
        pi.Va = msg.Va
        pi.Nh = msg.Nh
        // reply accept OK
        newPkt := pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = paxosproto.ACCEPT_OK
        pi.out <- newPkt
    }
}

func (pi *PaxosInstance) handleCommit() {
	//receive commit, notify paxosengine to record the log and take action
	newPkt := pi.initPkt()
	newPkt.PeerID = pi.PeerID
	newPkt.Msg.Type = paxosproto.COMMIT_OK
	newPkt.Msg.Va = pi.Va
	pi.prog <- newPkt
}

func (pi *PaxosInstance) Prepare() (int, * paxosproto.ValueStruct) {
    msg := &paxosproto.MsgStruct{}
    msg.Type = paxosproto.PREPARE
    // msg.Na = pi.Myn
    newPkt := &paxosproto.Packet{pi.PeerID,msg}
    pi.brd <- newPkt 
    p := <- pi.prepareCh //wait for the reply
    var state int
    var v * paxosproto.ValueStruct = nil
    switch p.Type {
    	case paxosproto.PREPARE_OK:
    		state = paxosproto.PREPARE_OK
    		v = &p.Va
    	case paxosproto.PREPARE_REJECT:
 			state = paxosproto.PREPARE_REJECT   	
    	case paxosproto.PREPARE_BEHIND:
    		state = paxosproto.PREPARE_BEHIND
    		v = &p.Va
    }
    return state, v
}

func (pi * PaxosInstance) Accept() int {
	msg := &paxosproto.MsgStruct{}
	msg.Type = paxosproto.ACCEPT
	msg.Myn = pi.Myn
	msg.Va = pi.Va
	newPkt := &paxosproto.Packet{pi.PeerID,msg}
	pi.brd <- newPkt
	p := <- pi.acceptCh
	switch p.Type {
		case paxosproto.ACCEPT_OK:
			return paxosproto.ACCEPT_OK	
		case paxosproto.ACCEPT_REJECT:
			return paxosproto.ACCEPT_REJECT
	}
	return -1	//error is -1
}

func (pi * PaxosInstance) Commit() {
	msg := &paxosproto.MsgStruct{}
	msg.Type = paxosproto.COMMIT
	msg.Va = pi.Va
	newPkt := &paxosproto.Packet{pi.PeerID,msg}
	pi.brd <- newPkt
}
