package paxosproto

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
	C_QueryFlights = iota
	C_PrepareBookFlight
	C_PrepareCancelFlight
	C_BookDecision
	C_CancelDecision
	C_DeleteFlight
	C_RescheduleFlight
	C_AddFlight
    C_NOP
)

const (
	Propose_OK = iota
	Propose_RETRY
	Propose_FAIL
)

type NodeStruct struct {
	Port string 
	ID int	
}

type ValueStruct struct {
    CoordSeq int
    Type int
    Action []byte
    Reply []byte
    Host string
}

type ReplyStruct struct {
	Reply []byte
	Type int
	Status int
}

type Packet struct {
    PeerID int
    Msg * MsgStruct
}

type MsgStruct struct {
    Seq int
    Type int
    Na int
    Myn int
    Va *ValueStruct
}

