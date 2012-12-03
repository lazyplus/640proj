package consensus

import (
	"net/rpc"
)

const (
	PREPARE = iota
	PREPARE_OK
	PREPARE_BEHIND
	PREPARE_REJECT
	ACCEPT
	ACCEPT_OK
	ACCEPT_REJECT
	COMMIT
)

const (
	c_QueryFlights = iota
	c_PrepareBookFlight
	c_PrepareCancelFlight
	c_BookDecision
	c_CancelDecision
	c_DeleteFlight
	c_RescheduleFlight
	c_AddFlight
)

type Delegate struct {
    Action interface{}
    timer chan interface{}
    servers []string
    cli [] *rpc.Client
    Port string
    numServers int
}

type NodeStruct struct {
	Port string 
	ID int	
}

type RPCStruct struct {
	in chan * Packet
}

type PaxosEngine struct {
	numNodes int
    log map[int] *ValueStruct
    clients []*rpc.Client
    servers []*NodeStruct
    cur_seq int
    cur_paxos PaxosInstance
    mutex sync.Mutex
    in chan * MsgStruct
    out chan * Packet
    brd chan * Packet
	prog chan * Packet
	exitCurrentPI chan int
    peerID int
    // network
    RPCReceiver * RPCStruct
}

type ValueStruct struct {
    CoordSeq int
    Type int
    Action interface{}
    Host string
}

type Packet struct {
    PeerID int
    Msg * MsgStruct
}

type PaxosInstance struct {
	log map[int] *ValueStruct
    out chan * Packet
    brd chan * Packet
    in chan * MsgStruct
    prog chan * Packet
    prepareCh chan * MsgStruct
    acceptCh chan * MsgStruct
    shouldExit chan int
	Na int
	Nh int
	Myn int
	Va ValueStruct
    seq int
    PeerID int
    numNodes int
    PreaccepteNodes int
    PrefailNodes int
    AcpacceptedNodes int
    AcpfailNodes int
}

type MsgStruct struct {
    Seq int
    Type int
    Na int
    Myn int
    Va ValueStruct
}

