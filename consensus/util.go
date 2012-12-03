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
    log map[int] *logStruct
    clients []*rpc.Client
    servers []*NodeStruct
    cur_seq int
    cur_paxos PaxosInstance
    mutex sync.Mutex
    in chan * Packet
    out chan * Packet
    brd chan * MsgStruct
    peerID int
    // network
    RPCReceiver * RPCStruct
}

type logStruct struct {
	op_type int
	Action interface{}
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

type PaxosInstance struct {
    out chan * Packet
    brd chan * Packet
    in chan * MsgStruct
    prepareCh chan * MsgStruct
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

