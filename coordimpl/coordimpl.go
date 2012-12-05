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
}

func NewCoordinator (path string) *coordserver {
	co := &coordserver{}
	co.coordSeq = 0
	co.airline_info , _ = config.ReadConfigFile(path)
	return co
}

func (co *coordserver) BookFlights(args *coordproto.BookArgs, ori_reply *coordproto.BookReply) error {
	//2pc
	//assume the flight ID is in format: airline + "::" + ID 
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	//first round
	ls := len(args.Flights)
	id_ls := list.New()
	wait_ls := make(map[string] *rpc.Call)
	reply_ls := make(map[string] *delegateproto.BookReply)
	client_ls := make(map[string] *rpc.Client)
	var shouldAbort bool = false
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
		args_out := &delegateproto.BookArgs{id,args.Email,args.Count,co.coordSeq}
		reply := &delegateproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		bookcall := client.Go("AirlineServerRPC.PrepareBookFlight",args_out, reply,nil)
		id_ls.PushBack(id)
		wait_ls[id] = bookcall
		reply_ls[id] = reply
		client_ls[id] = client	
	}
	var shouldCommit int = delegateproto.COMMIT
	var finalstatus int = coordproto.OK
	
	for e := id_ls.Front();e != nil; e=e.Next() {
		ss := e.Value.(string)
		call , _:= wait_ls[ss]
		<- call.Done
		reply , _:= reply_ls[ss]
		if reply.Status != delegateproto.OK || shouldAbort == true {
			shouldCommit = delegateproto.ABORT
			finalstatus = reply.Status
		}
	}
	
	//second phase
	wait_ls2 := make(map[string] *rpc.Call)
	reply_ls2 := make(map[string] *delegateproto.DecisionReply)
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		client, _ := client_ls[ss]
		args_out := &delegateproto.DecisionArgs{shouldCommit,ss,co.coordSeq}
		reply := &delegateproto.DecisionReply{}
		decisioncall := client.Go("DelegateServerRPC.BookDecision", args_out, reply,nil)
		wait_ls2[ss] = decisioncall
		reply_ls2[ss] = reply
	}
	
	
	for e:=id_ls.Front();e!=nil;e=e.Next(){
		ss := e.Value.(string)
		call , _ := wait_ls2[ss]
		<- call.Done
		reply , _ := reply_ls2[ss]
		if reply.Status != delegateproto.OK {
			finalstatus = reply.Status
		}
		cl := client_ls[ss]
		cl.Close()
	}	
	ori_reply.Status = finalstatus
	if shouldAbort {
		ori_reply.Status = coordproto.ENOFLIGHT
	}
	co.coordSeq ++ 
	return nil
}

func (co *coordserver) CancelFlights(args *coordproto.BookArgs, ori_reply *coordproto.BookReply) error {
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	
	ls := len(args.Flights)
	id_ls := list.New()
	wait_ls := make(map[string] *rpc.Call)
	reply_ls := make(map[string] *delegateproto.BookReply)
	client_ls := make(map[string] *rpc.Client)
	//assume the input is airline + "::" + ID
	var shouldAbort bool = false
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
		args_out := &delegateproto.BookArgs{id,args.Email,args.Count,co.coordSeq}
		reply := &delegateproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		cancelcall := client.Go("DelegateServerRPC.PrepareCancelFlight",args_out,reply,nil)
		id_ls.PushBack(id)
		client_ls[id] = client
		wait_ls[id] = cancelcall
		reply_ls[id] = reply
	}	
	var should_commit int = delegateproto.COMMIT
	var final_status int = coordproto.OK
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		call, _ := wait_ls[ss]
		<- call.Done
		reply,_ := reply_ls[ss]
		if reply.Status != delegateproto.OK || shouldAbort {
			should_commit = delegateproto.ABORT
			final_status = reply.Status
		}
	}
	//second phase
	wait_ls2 := make(map[string] *rpc.Call)
	reply_ls2 := make(map[string] *delegateproto.DecisionReply)

	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		client, _ := client_ls[ss]
		args_out := &delegateproto.DecisionArgs{should_commit,ss,co.coordSeq}
		reply := &delegateproto.DecisionReply{}
		call := client.Go("AirlineServerRPC.CancelDecision",args_out,reply, nil)
		wait_ls2[ss] = call
		reply_ls2[ss] = reply
	}
	
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		call, _ := wait_ls2[ss]
		<- call.Done
		reply, _ := reply_ls2[ss]
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
	co.coordSeq++
	return nil	
}


func (co *coordserver) QueryFlights(args * coordproto.QueryArgs, reply * coordproto.QueryReply) error {
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	return nil
}

func (co *coordserver) DeleteFlight(* coordproto.DeleteArgs, * coordproto.DeleteReply) error {
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	return nil
}

func (co *coordserver) RescheduleFlight(* coordproto.RescheduleArgs, * coordproto.RescheduleReply) error {
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	return nil
}

func (co *coordserver) AddFlight(args * coordproto.AddArgs, reply * coordproto.AddReply) error {
	co.map_lock.Lock()
	defer co.map_lock.Unlock()

	fmt.Println("Adding Flight")
	fmt.Println(args.Flight)
	dargs := &delegateproto.AddArgs{}
	dargs.Flight = args.Flight
	dargs.Seqnum = co.coordSeq
	co.coordSeq ++
	asname := strings.Split(dargs.Flight.FlightID, "-")[0]
	fmt.Println("Adding to airline " + asname)
	var dreply delegateproto.AddReply
	cli, err := rpc.DialHTTP("tcp", co.airline_info.Airlines[asname].DelegateHostPort)
	if err == nil {
		cli.Call("DelegateServerRPC.AddFlight", dargs, &dreply)
		reply.Status = dreply.Status
		return nil
	}
	reply.Status = coordproto.ENOAIRLINE
	return nil
}
