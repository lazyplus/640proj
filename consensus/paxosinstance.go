package consensus

type PaxosInstance struct {
    out * chan Packet
    brd * chan Packet
    in * chan Msg

    seq int
}

type MsgStruct struct {
    Seq int
    Type int
    Na int
    Va ValueStruct
}

func NewPaxosInstance(path string) *PaxosInstance {
    pi := &PaxosInstance{}
    // GetConfig
}

func (pi *PaxosInstance) Run() {
    for {
        select{
        case inPkt := <- pi.in:
            switch(inPkt.Msg.Type){
            case PREPARE:
                handlePrepare(inPkt)
            case PREPARE_OK:
                fallthrough
            case PREPARE_REJECT:
                fallthrough
            case PREPARE_BEHIND:
                pi.prepareCh <- inPkt.Msg.Type
            case ACCEPT:
                handleAccept(inPkt)
            case ACCEPT_OK:
                fallthrough
            case ACCEPT_REJECT:
                pi.acceptCh <- inPkt.Msg.Type
            case COMMIT:
                handleCommit(inPkt)
            }
        }
    }
}

func (pi *PaxosInstance) handlePrepare(pkt Packet) {
    if pkt.Msg.Seq < pi.seq {
        // reply prepare behind
        newPkt = &Packet{}
        newPkt.Msg = &Msg{}
        newPkt.PeerID = pkt.PeerID
        newPkt.Msg.Seq = pi.seq
        newPkt.Msg.Type = PREPARE_BEHIND
        newPkt.Msg.Va = 
    }
    if msg.Na < pi.Na {
        // reply prepare rejected
    }else {
        pi.Nh = msg.Na
        // reply prepare OK, Va
    }
}

func (pi *PaxosInstance) handleAccept(msg Msg) {
    if msg.Na < pi.Nh {
        // reply accept rejected
    }else{
        pi.Na = msg.Na
        pi.Va = msg.Va
        pi.Nh = msg.Nh
        // reply accept OK
    }
}

func (pi *PaxosInstance) handleCommit() {

}

func (pi *PaxosInstance) prepare() {
    msg := &MsgStruct{}
    msg.Type = PREPARE
    msg.Na = pi.Myn
    pi.brd <- msg
    return <- pi.prepareCh
}
