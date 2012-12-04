package consensus

import (
	"math/rand"
)

func generate_random_number(PeerID int, numNodes int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return (r.Intn(10)*numNodes + peerID)
}

func NewPaxosInstance(peerID int, seq int, numNodes int) *PaxosInstance {
    pi := &PaxosInstance{}
    pi.seq = seq
    pi.PeerID = peerID
    pi.numNodes = numNodes
    pi.Myn = generate_random_number(peerID,numNodes)
    pi.prepareCh = make(chan * MsgStruct)
    pi.acceptCh = make(chan * MsgStruct)
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
            if inPkt.Msg.Type == PREPARE {
                handlePrepare(inPkt)
                break
            }

            if inPkt.Msg.Seq != pi.seq {
                break
            }

            switch(inPkt.Msg.Type){
            case PREPARE_OK:
				pi.PreaccepteNodes ++
                if pi.Na < inPkt.Msg.Na {
                    pi.Na = inPkt.Msg.Na
                    pi.Va = inPkt.Msg.Va
                }
				if (pi.PreaccepteNodes > (numNodes/2)) && !finishPrepare {
					finishPrepare = true
					pi.prepareCh <- inPkt.Msg
				}
            case PREPARE_REJECT:
                pi.PrefailNodes ++
                if pi.PrefailNodes > (numNodes/2)) && !finishPrepare {
                	finishPrepare = true
 					pi.prepareCh <- inPkt.Msg
                }
            case PREPARE_BEHIND:
                pi.Va = inPkt.Msg.Va
            	if !finishPrepare {
            		finishPrepare = true
					pi.prepareCh <- inPkt.Msg
                }
            case ACCEPT:
                handleAccept(inPkt)
            case ACCEPT_OK:
                pi.AcpacceptedNodes ++
                if (pi.AcpacceptedNodes > (numNodes/2)) && !finishAccept{
                	finishAccept = true
                	pi.acceptCh <- inPkt.Msg
                }
            case ACCEPT_REJECT:
				pi.AcpfailNodes ++
				if (pi.AcpacceptedNodes > (numNodes/2)) && !finishAccept{
                	finishAccept = true
                	pi.acceptCh <- inPkt.Msg
                }
            case COMMIT:
                handleCommit(inPkt)
            }
        }
    }
}

func (pi *PaxosInstance) initPkt () *Packet {
	p := &Packet{}
	p.Msg = &MsgStruct{}
	p.PeerID = pi.PeerID
	p.Msg.Seq = pi.seq
	return p
}

func (pi *PaxosInstance) handlePrepare(pkt Packet) {
    if pkt.Msg.Seq < pi.seq {
        // reply prepare behind
        newPkt = pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = PREPARE_BEHIND
        newPkt.Msg.Va = pi.log[pkt.Msg.Seq]
        pi.out <- newPkt
    }
    if msg.Na < pi.Nh {
        // reply prepare rejected
 		newPkt = pi.initPkt()
 		newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = PREPARE_REJECT
        pi.out <- newPkt
    }else {
        pi.Nh = msg.Na
        // reply prepare OK, Va
        newPkt = pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = PREPARE_OK
        newPkt.Msg.Va = pi.Va
        newPkt.Msg.Na = pi.Na
        pi.out <- newPkt
    }
}

func (pi *PaxosInstance) handleAccept(pkt Packet) {
	msg := pkt.Msg
    if msg.Na < pi.Nh {
        // reply accept rejected
        newPkt = pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = ACCEPT_REJECT
        pi.out <- newPkt
    }else{
        pi.Na = msg.Na
        pi.Va = msg.Va
        pi.Nh = msg.Nh
        // reply accept OK
        newPkt = pi.initPkt()
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Type = ACCEPT_OK
        pi.out <- newPkt
    }
}

func (pi *PaxosInstance) handleCommit() {
	//receive commit, notify paxosengine to record the log and take action
	newPkt = pi.initPkt()
	newPkt.PeerID = pi.PeerID
	newPkt.Msg.Type = COMMIT_OK
	newPkt.Msg.Va = pi.Va
	pi.prog <- newPkt
}

func (pi *PaxosInstance) Prepare() (int, ValueStruct) {
    msg := &MsgStruct{}
    msg.Type = PREPARE
    // msg.Na = pi.Myn
    newPkt := &Packet{pi.PeerID,msg}
    pi.brd <- newPkt 
    p := <- pi.prepareCh //wait for the reply
    var state int
    var v ValueStruct
    switch p.Type {
    	case PREPARE_OK:
    		state = PREPARE_OK
    		v = p.Va
    	case PREPARE_REJECT:
 			state = PREPARE_REJECT   	
    	case PREPARE_BEHIND:
    		state = PREPARE_BEHIND
    		v = p.Va
    }
    return state, v
}

func (pi * PaxosInstance) Accept() int {
	msg := &MsgStruct{}
	msg.Type = ACCEPT
	msg.Myn = pi.Myn
	msg.Va = pi.Va
	newPkt := &Packet{pi.PeerID,msg}
	pi.brd <- newPkt
	p := <- pi.acceptCh
	switch p.Type {
		case ACCEPT_OK:
			return ACCEPT_OK	
		case ACCEPT_REJECT:
			return ACCEPT_REJECT
	}
	return -1	//error is -1
}

func (pi * PaxosInstance) Commit() {
	msg := &MsgStruct{}
	msg.Type = COMMIT
	msg.Va = pi.Va
	newPkt := &Packet{pi.PeerID,msg}
	pi.brd <- newPkt
}
