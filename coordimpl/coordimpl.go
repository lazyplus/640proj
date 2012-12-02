package coordimpl

import (
	"../coordproto"
	"../airlineproto"
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
}

func NewCoordinator (path string) *coordserver {
	co := &coordserver{}
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
	reply_ls := make(map[string] *airlineproto.BookReply)
	client_ls := make(map[string] *rpc.Client)
	var shouldAbort bool = false
	for i:=0;i<ls;i++ {
		fmt.Println("Flight: " + args.Flights[i])
		ss := strings.Split(args.Flights[i],"-")
		id := args.Flights[i]
		airline_name := ss[0]
		addr , found := co.airline_info.AirlineAddr[airline_name]
		if found == false {
			shouldAbort = true
			break
		}
		args_out := &airlineproto.BookArgs{id,args.Email,args.Count}
		reply := &airlineproto.BookReply{}
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
	var shouldCommit int = airlineproto.COMMIT
	var finalstatus int = coordproto.OK
	
	for e := id_ls.Front();e != nil; e=e.Next() {
		ss := e.Value.(string)
		call , _:= wait_ls[ss]
		<- call.Done
		reply , _:= reply_ls[ss]
		if reply.Status != airlineproto.OK || shouldAbort == true {
			shouldCommit = airlineproto.ABORT
			finalstatus = reply.Status
		}
	}
	
	//second phase
	wait_ls2 := make(map[string] *rpc.Call)
	reply_ls2 := make(map[string] *airlineproto.DecisionReply)
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		client, _ := client_ls[ss]
		args_out := &airlineproto.DecisionArgs{shouldCommit,ss}
		reply := &airlineproto.DecisionReply{}
		decisioncall := client.Go("AirlineServerRPC.BookDecision", args_out, reply,nil)
		wait_ls2[ss] = decisioncall
		reply_ls2[ss] = reply
	}
	
	
	for e:=id_ls.Front();e!=nil;e=e.Next(){
		ss := e.Value.(string)
		call , _ := wait_ls2[ss]
		<- call.Done
		reply , _ := reply_ls2[ss]
		if reply.Status != airlineproto.OK {
			finalstatus = reply.Status
		}
		cl := client_ls[ss]
		cl.Close()
	}	
	ori_reply.Status = finalstatus
	if shouldAbort {
		ori_reply.Status = coordproto.ENOFLIGHT
	}
	return nil
}

func (co *coordserver) CancelFlights(args *coordproto.BookArgs, ori_reply *coordproto.BookReply) error {
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	
	ls := len(args.Flights)
	id_ls := list.New()
	wait_ls := make(map[string] *rpc.Call)
	reply_ls := make(map[string] *airlineproto.BookReply)
	client_ls := make(map[string] *rpc.Client)
	//assume the input is airline + "::" + ID
	var shouldAbort bool = false
	for i:=0;i<ls;i++ {
		ss := strings.Split(args.Flights[i],"-")
		airline_name := ss[0]
		id := args.Flights[i]
		addr, found := co.airline_info.AirlineAddr[airline_name]
		if found == false{
			shouldAbort = true
			break
		}
		args_out := &airlineproto.BookArgs{id,args.Email,args.Count}
		reply := &airlineproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		cancelcall := client.Go("AirlineServerRPC.PrepareCancelFlight",args_out,reply,nil)
		id_ls.PushBack(id)
		client_ls[id] = client
		wait_ls[id] = cancelcall
		reply_ls[id] = reply
	}	
	var should_commit int = airlineproto.COMMIT
	var final_status int = coordproto.OK
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		call, _ := wait_ls[ss]
		<- call.Done
		reply,_ := reply_ls[ss]
		if reply.Status != airlineproto.OK || shouldAbort {
			should_commit = airlineproto.ABORT
			final_status = reply.Status
		}
	}
	//second phase
	wait_ls2 := make(map[string] *rpc.Call)
	reply_ls2 := make(map[string] *airlineproto.DecisionReply)

	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		client, _ := client_ls[ss]
		args_out := &airlineproto.DecisionArgs{should_commit,ss}
		reply := &airlineproto.DecisionReply{}
		call := client.Go("AirlineServerRPC.CancelDecision",args_out,reply, nil)
		wait_ls2[ss] = call
		reply_ls2[ss] = reply
	}
	
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.Value.(string)
		call, _ := wait_ls2[ss]
		<- call.Done
		reply, _ := reply_ls2[ss]
		if reply.Status != airlineproto.OK {
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
