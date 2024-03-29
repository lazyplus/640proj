package coordimpl

import (
	"../coordproto"
	"../delegateproto"
	"../config"
	"sync"
	"strings"
	"net/rpc"
	"container/list"
	"fmt"
)

type coordserver struct {
	airline_info	*config.Config
	map_lock		sync.Mutex
	coordSeq		int
	work_slot chan int
}

func NewCoordinator (path string) *coordserver {
	co := &coordserver{}
	co.coordSeq = 0
	co.work_slot = make(chan int, 5)
	for i:=0; i<5; i++ {
		co.work_slot <- 1
	}
	co.airline_info , _ = config.ReadConfigFile(path)
	return co
}

func (co *coordserver) BookFlights(args *coordproto.BookArgs, ori_reply *coordproto.BookReply) error {
	<- co.work_slot
	defer co.returnSlot()
	//2pc
	//assume the flight ID is in format: airline + "-" + ID 

	co.map_lock.Lock()
	co.coordSeq ++ 
	this_coordseq := co.coordSeq
	co.map_lock.Unlock()
	
	//first round
	ls := len(args.Flights)
	id_ls := list.New()
	client_ls := make(map[string] *rpc.Client)
	var shouldAbort bool = false
	var shouldCommit int = delegateproto.COMMIT
	var finalstatus int = coordproto.OK
	
	fmt.Println("Booking")
	for i:=0;i<ls;i++ {
		fmt.Println("Flight: " + args.Flights[i])
		ss := strings.Split(args.Flights[i],"-")
		id := args.Flights[i]
		airline_name := ss[0]
		addr_list , found := co.airline_info.Airlines[airline_name]
		if found == false {
			shouldAbort = true
			break
		}
		addr := addr_list.DelegateHostPort
		args_out := &delegateproto.BookArgs{id,args.Email,args.Count,this_coordseq}
		reply := &delegateproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		client.Call("DelegateServerRPC.PrepareBookFlight",args_out, reply)
		if reply.Status != delegateproto.OK || shouldAbort == true {
			shouldCommit = delegateproto.ABORT
			finalstatus = reply.Status
			shouldAbort = true
		}
		id_ls.PushBack(id)
		client_ls[id] = client	
	}
	
	co.map_lock.Lock()
	co.coordSeq ++ 
	this_coordseq = co.coordSeq
	co.map_lock.Unlock()
	//second phase

	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		client, _ := client_ls[ss]
		args_out := &delegateproto.DecisionArgs{shouldCommit,ss,this_coordseq}
		reply := &delegateproto.DecisionReply{}
		client.Call("DelegateServerRPC.BookDecision", args_out, reply)
		if reply.Status != delegateproto.OK {
			finalstatus = reply.Status
		}
		cl := client_ls[ss]
		cl.Close()
	}
	
	ori_reply.Status = finalstatus

	return nil
}

func (co *coordserver) CancelFlights(args *coordproto.BookArgs, ori_reply *coordproto.BookReply) error {
	<- co.work_slot
	defer co.returnSlot()

	co.map_lock.Lock()
	co.coordSeq ++ 
	this_coordseq := co.coordSeq
	co.map_lock.Unlock()

	ls := len(args.Flights)
	id_ls := list.New()

	client_ls := make(map[string] *rpc.Client)
	//assume the input is airline + "-" + ID
	var shouldAbort bool = false
	var should_commit int = delegateproto.COMMIT
	var final_status int = coordproto.OK
	for i:=0;i<ls;i++ {
		ss := strings.Split(args.Flights[i],"-")
		airline_name := ss[0]
		id := args.Flights[i]
		addr_list, found := co.airline_info.Airlines[airline_name]
		if found == false{
			shouldAbort = true
			break
		}
		addr := addr_list.DelegateHostPort
		args_out := &delegateproto.BookArgs{id,args.Email,args.Count,this_coordseq}
		reply := &delegateproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		client.Call("DelegateServerRPC.PrepareCancelFlight",args_out,reply)
		id_ls.PushBack(id)
		client_ls[id] = client
		if reply.Status != delegateproto.OK || shouldAbort {
			should_commit = delegateproto.ABORT
			final_status = reply.Status
		}
	}

	co.map_lock.Lock()
	co.coordSeq ++ 
	this_coordseq = co.coordSeq
	co.map_lock.Unlock()
	//second phase

	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		client, _ := client_ls[ss]
		args_out := &delegateproto.DecisionArgs{should_commit,ss,this_coordseq}
		reply := &delegateproto.DecisionReply{}
		client.Call("DelegateServerRPC.CancelDecision",args_out,reply)
		if reply.Status != delegateproto.OK {
			final_status = reply.Status
		}
		cl := client_ls[ss]
		cl.Close()
	}
	
	ori_reply.Status = final_status
	if shouldAbort {
		ori_reply.Status = coordproto.ENOFLIGHT
	}
	return nil	
}


func (co *coordserver) QueryFlights(args * coordproto.QueryArgs, reply * coordproto.QueryReply) error {
	<- co.work_slot
	defer co.returnSlot()

	co.map_lock.Lock()
	co.coordSeq ++
	this_coordseq := co.coordSeq
	co.map_lock.Unlock()

	reply.FlightList = make([]delegateproto.FlightStruct, 0)

	for _, value := range(co.airline_info.Airlines) {
		dargs := &delegateproto.QueryArgs{}
		dargs.Seqnum = this_coordseq
		dargs.StartTime = args.StartTime
		dargs.EndTime = args.EndTime
		// fmt.Println("Querying airline " + key)
		var dreply delegateproto.QueryReply
		cli, err := rpc.DialHTTP("tcp", value.DelegateHostPort)
		if err == nil {
			cli.Call("DelegateServerRPC.QueryFlights", dargs, &dreply)
			reply.Status = dreply.Status
			reply.FlightList = append(reply.FlightList, dreply.FlightList...)
			cli.Close()
		}
	}
	return nil
}

func (co *coordserver) DeleteFlight(args * coordproto.DeleteArgs, reply * coordproto.DeleteReply) error {
	<- co.work_slot
	defer co.returnSlot()

	dargs := &delegateproto.DeleteArgs{}

	dargs.FlightID = args.FlightID
	co.map_lock.Lock()
	co.coordSeq ++
	dargs.Seqnum = co.coordSeq
	co.map_lock.Unlock()
	
	asname := strings.Split(dargs.FlightID, "-")[0]
	var dreply delegateproto.DeleteReply
	cli, err := rpc.DialHTTP("tcp", co.airline_info.Airlines[asname].DelegateHostPort)
	if err == nil {
		cli.Call("DelegateServerRPC.DeleteFlight", dargs, &dreply)
		reply.Status = dreply.Status
		reply.CustomerEmails = dreply.CustomerEmails
		cli.Close()
	}
	return nil
}

func (co *coordserver) RescheduleFlight(args * coordproto.RescheduleArgs, reply * coordproto.RescheduleReply) error {
	<- co.work_slot
	defer co.returnSlot()

	dargs := &delegateproto.RescheduleArgs{}
	dargs.OldFlightID = args.OldFlightID
	dargs.NewFlight = args.NewFlight

	co.map_lock.Lock()
	co.coordSeq ++
	dargs.Seqnum = co.coordSeq
	co.map_lock.Unlock()

	asname := strings.Split(dargs.OldFlightID, "-")[0]
	var dreply delegateproto.RescheduleReply
	cli, err := rpc.DialHTTP("tcp", co.airline_info.Airlines[asname].DelegateHostPort)
	if err == nil {
		cli.Call("DelegateServerRPC.RescheduleFlight", dargs, &dreply)
		reply.Status = dreply.Status
		reply.CustomerEmails = dreply.CustomerEmails
		cli.Close()
		return nil
	}
	reply.Status = coordproto.ENOAIRLINE
	return nil
}

func (co *coordserver) AddFlight(args * coordproto.AddArgs, reply * coordproto.AddReply) error {
	<- co.work_slot
	defer co.returnSlot()

	dargs := &delegateproto.AddArgs{}
	dargs.Flight = args.Flight

	co.map_lock.Lock()
	co.coordSeq ++
	dargs.Seqnum = co.coordSeq
	co.map_lock.Unlock()

	asname := strings.Split(dargs.Flight.FlightID, "-")[0]
	var dreply delegateproto.AddReply
	cli, err := rpc.DialHTTP("tcp", co.airline_info.Airlines[asname].DelegateHostPort)
	if err == nil {
		cli.Call("DelegateServerRPC.AddFlight", dargs, &dreply)
		reply.Status = dreply.Status
		cli.Close()
		return nil
	}
	reply.Status = coordproto.ENOAIRLINE
	return nil
}

func (co *coordserver) returnSlot() {
	co.work_slot <- 1
}
