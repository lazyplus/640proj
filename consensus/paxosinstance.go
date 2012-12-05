package consensus

import (
	// "math/rand"
	"../paxosproto"
	// "time"
    // "fmt"
)

func generate_random_number(PeerID int, numNodes int, Nh int) int {
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
    // fmt.Println("Generating MyN")
    // fmt.Println((Nh / 10 + 1) * 10 + PeerID)
    return (Nh / 10 + 1) * 10 + PeerID
	// return (r.Intn(10)*numNodes + PeerID)
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
	// finishPrepare := false
	// finishAccept := false
	// finishCommit := false
	
    for pi.running {
        select{
        case <- (*pi.shouldExit):
        	pi.running = false
        	break
        case inPkt := <- pi.in:
            // fmt.Println("Received ")
            // fmt.Println(inPkt)
			if pi.running == false {
				break
			}
            if inPkt.Msg.Type == paxosproto.PREPARE {
                pi.handlePrepare(inPkt)
                break
            }

            if inPkt.Msg.Seq < pi.seq {
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
                pi.handleAccept(inPkt)
            case paxosproto.ACCEPT_OK:
                // fmt.Println("Get ACCEPT_OK")
                pi.AcpacceptedNodes ++
                if (pi.AcpacceptedNodes > (pi.numNodes/2)) {
                	pi.acceptCh <- inPkt.Msg
                }
            case paxosproto.ACCEPT_REJECT:
                // fmt.Println("Get ACCEPT_REJECT")
				pi.AcpfailNodes ++
				if (pi.AcpacceptedNodes > (pi.numNodes/2)){
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

func (pi *PaxosInstance) handleAccept(pkt * paxosproto.Packet) {
    // fmt.Println("Handling Accept")
    // fmt.Println(pkt)
	msg := pkt.Msg
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
    pi.Finished = true
	pi.prog <- pkt.Msg.Va
}

func (pi *PaxosInstance) Prepare() (int, * paxosproto.ValueStruct) {
    pi.Myn = generate_random_number(pi.PeerID,pi.numNodes, pi.Nh)
    msg := &paxosproto.MsgStruct{}
    msg.Type = paxosproto.PREPARE
    msg.Seq = pi.seq
    msg.Na = pi.Myn
    newPkt := &paxosproto.Packet{pi.PeerID,msg}
    // fmt.Println("Sending prepare msg")
    // fmt.Println(newPkt.Msg)
    pi.brd <- newPkt 
    // fmt.Println("Waiting for prepare ok")
    p := <- pi.prepareCh //wait for the reply
    // fmt.Println("prepare done")
    // fmt.Println(p)
    var state int
    var v * paxosproto.ValueStruct = nil
    switch p.Type {
    	case paxosproto.PREPARE_OK:
    		state = paxosproto.PREPARE_OK
    		v = p.Va
    	case paxosproto.PREPARE_REJECT:
 			state = paxosproto.PREPARE_REJECT   	
    	case paxosproto.PREPARE_BEHIND:
    		state = paxosproto.PREPARE_BEHIND
    		v = p.Va
    }
    // fmt.Println("Returning ")
    // fmt.Println(state)
    // fmt.Println(v)
    return state, v
}

func (pi * PaxosInstance) Accept(v *paxosproto.ValueStruct) int {
	msg := &paxosproto.MsgStruct{}
	msg.Type = paxosproto.ACCEPT
	msg.Na = pi.Myn
	msg.Va = v
    msg.Seq = pi.seq
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
    msg.Seq = pi.seq
	newPkt := &paxosproto.Packet{pi.PeerID,msg}
	pi.brd <- newPkt
}
