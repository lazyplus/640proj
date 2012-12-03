package consensus

import (
	"math/rand"
)

func generate_random_number(int PeerID) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return (r.Intn(10)*3 + peerID)
}

func NewPaxosInstance(peerID int, seq int, numNodes int) *PaxosInstance {
    pi := &PaxosInstance{}
    pi.seq = seq
    pi.PeerID = peerID
    pi.numNodes = numNodes
    pi.Myn = generate_random_number(peerID)
    return pi
}

func (pi *PaxosInstance) Run() {
	finishPrepare := false
	finishAccept := false
	finishCommit := false
	
    for {
        select{
        case inPkt := <- pi.in:
            switch(inPkt.Msg.Type){
            case PREPARE:
                handlePrepare(inPkt)
            case PREPARE_OK:
				pi.PreaccepteNodes ++
				if (pi.PreaccepteNodes > (numNodes/2)) && !finishPrepare {
					finishPrepare = true
//					pi.prepareCh <- inPkt.Msg.Type
					pi.prepareCh <- inPkt.Msg
				}                
            case PREPARE_REJECT:
                pi.PrefailNodes ++
                if pi.PrefailNodes > (numNodes/2)) && !finishPrepare {
                	finishPrepare = true
//                	pi.prepareCh <- inPkt.Msg.Type
 					pi.prepareCh <- inPkt.Msg
                }
            case PREPARE_BEHIND:
            	if !finishPrepare {
            		finishPrepare = true
//                	pi.prepareCh <- inPkt.Msg.Type
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
//                pi.acceptCh <- inPkt.Msg.Type
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
	p.Msg = &MsgStruct
	p.PeerID = pi.PeerID
	p.Msg.Seq = pi.seq
	return p
}

func (pi *PaxosInstance) handlePrepare(pkt Packet) {
    if pkt.Msg.Seq < pi.seq {
        // reply prepare behind
        newPkt = pi.initPkt()
        newPkt.Msg.Type = PREPARE_BEHIND
        out <- newPkt
    }
    if msg.Na < pi.Na {
        // reply prepare rejected
 		newPkt = pi.initPkt()
        newPkt.Msg.Type = PREPARE_REJECT
        out <- newPkt
    }else {
        pi.Nh = msg.Na
        // reply prepare OK, Va
        newPkt = pi.initPkt()
        newPkt.Msg.Type = PREPARE_OK
        newPkt.Msg.Va = pi.Va
        newPkt.Msg.Na = pi.Na
        out <- newPkt
    }
}


func (pi *PaxosInstance) handleAccept(msg Msg) {
    if msg.Na < pi.Nh {
        // reply accept rejected
        newPkt = pi.initPkt()
        newPkt.Msg.Type = ACCEPT_REJECT
        out <- newPkt
    }else{
        pi.Na = msg.Na
        pi.Va = msg.Va
        pi.Nh = msg.Nh
        // reply accept OK
        newPkt.Msg.Type = ACCEPT_OK
        out <- newPkt
    }
}

func (pi *PaxosInstance) handleCommit() {
	newPkt = pi.initPkt()
	newPkt.Msg.Type = COMMIT_OK
	out <- newPkt
}

func (pi *PaxosInstance) prepare() (int, ValueStruct) {
    msg := &MsgStruct{}
    msg.Type = PREPARE
    msg.Na = pi.Myn
    pi.brd <- msg 
    p := <- pi.prepareCh //wait for the reply
    var state int
    var v ValueStruct
    switch p.Type {
    	case PREPARE_OK:
    		state = PREPARE_OK
    		v = p.Va
    	case PREPARE_REJECT:
 			state = PREPARE_REJECT
 			v = nil   	
    	case PREPARE_BEHIND:
    		state = PREPARE_BEHIND
    		v = nil
    }
    return state, v
}

func (pi * PaxosInstance) Accept() int {
	msg := &MsgStruct{}
	msg.Type = ACCEPT
	msg.Myn = pi.Myn
	msg.Va = pi.Va
	pi.brd <- msg
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
	pi.brd <- msg
}
