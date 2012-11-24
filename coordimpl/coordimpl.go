package coordimpl

import {
	"../coordproto"
	"../airlineproto"
	"../inputFormat"
	"sync"
	"os"
	"bufio"
	"strings"
	"net/rpc"
	"container/list"
}

type coordserver struct {
	airline_info	*inputFormat.input
	map_lock		sync.Mutex
}

func NewCoordinator (path string) *coordserver {
	co := &coordserver{}
	co.airline_info = *inputFormat.input{}
	co.airline_info.read_config_file(path)
	return co
}

func (co *coordserver) BookFlights(args *coordproto.BookArgs, ori_reply *coordproto.BookReply) error {
	//2pc
	//assume the flight ID is in format: airline + "::" + ID 
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	//first round
	ls := len(args.Flights)
	id_ls := List.new()
	wait_ls := make(map[string] *rpc.Call)
	reply_ls := make(map[string] *airlineproto.BookReply)
	client_ls := make(map[string] *rpc.Client)
	var shouldAbort bool = false
	for i:=0;i<ls;i++ {
		ss := strings.Split(args.Flights[i],"-")
		airline_name = ss[0]
		addr , found := co.airline_info.addr_airline[airline_name]
		if found == false {
			shouldAbort = true
			break
		}
		args_out := &airlineproto.BookArgs{ss[1],args.Email,args.count}
		reply := &airlineproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		bookcall := client.Go("asrpc.PrepareBookFlight",args_out, reply,nil)
		id_ls.PushBack(ss)
		wait_ls[ss] = bookcall
		reply_ls[ss] = reply
		client_ls[ss] = client	
	}
	var shouldCommit int = airlineproto.COMMIT
	var finalstatus int = coordserver.OK
	
	for e := id_ls.Front();e != nil; e=e.Next() {
		ss := e.(string)
		call , _:= wait_ls[ss]
		reply_call := <- call.Done
		reply , _:= reply_ls[ss]
		if reply.Status != airlineproto.OK || shouldAbort == true {
			shouldCommit = airlineproto.ABORT
			finalstatus = reply.Status
		}
	}
	
	//second phase
	wait_ls = make(map[string] *rpc.Call)
	reply_ls = make(map[string] *airlineproto.DecisionReply)
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.(string)
		client, _ := client_ls[ss]
		args_out := &airlineproto.DecisionArgs{shouldCommit,ss}
		reply := &airlineproto.DecisionReply{}
		decisioncall := client.Go("asrpc.BookDecision", args_out, reply,nil)
		wait_ls[ss] = decisioncall
		reply_ls[ss] = reply
	}
	
	var isSuccess bool = true
	
	for e:=id_ls.Front();e!=nil;e=e.Next(){
		ss := e.(string)
		call , _ := wait_ls[ss]
		reply_call := <- call.Done
		reply , _ := reply_ls[ss]
		if reply.Status != airlineproto.OK {
			isSuccess = false
			finalstatus = reply.Status
		}
	}	
	ori_reply.Status = finalstatus
	if shouldAbort {
		ori_reply.Status = coordproto.ENOFLIGHT
	}
	return nil
}

func (co *coordserver) CancelFlights(args *coordproto.CancelArgs, ori_reply *coordproto.CancelReply) error {
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	
	ls := len(args.Flights)
	id_ls := List.new()
	wait_ls := make(map[string] *rpc.Call)
	reply_ls := make(map[string] *airlineproto.BookReply)
	client_ls := make(map[string] *rpc.Client)
	//assume the input is airline + "::" + ID
	var shouldAbort bool = false
	for i:=0;i<ls;i++ {
		ss := strings.Split(args.Flights[i],"-")
		airline_name := ss[0]
		airline_id := ss[1]
		addr, found := co.airline_info.addr_airline[airline_name]
		if found == false{
			shouldAbort = true
			break
		}
		args_out := &airlineproto.BookArgs{airline_id,args.Email,args.count}
		reply := &airlineproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		cancelcall := client.Go("asrpc.repareCancelFlight",args_out,reply,nil)
		id_ls.PushBack(ss)
		client_ls[ss] = client
		wait_ls[ss] = cancelcall
		reply_ls[ss] = reply
	}	
	var should_commit int = airlineproto.COMMIT
	var final_status int = coordproto.OK
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.(string)
		call, _ := wait_ls[ss]
		reply_call := call.Done
		reply,_ := reply_ls[ss]
		if reply.Status != airlineproto.OK || shouldAbort {
			should_commit = airlineproto.ABORT
			final_status = reply.Status
		}
	}
	//second phase
	wait_ls = make(map[string] *rpc.Call)
	reply_ls = make(map[string] *airlineproto.DecisionReply)
	var isSuccess bool = true
	for e:=id_ls.Front();e!=nil;e=e.Next {
		ss := e.(string)
		client, _ := client_ls[ss]
		args_out := &airlineproto.DecisionArgs{should_commit,ss}
		reply := &airlineproto.DecisionReply{}
		call := client.Go("asrpc.CancelDecision",args_out,reply)
		wait_ls[ss] = call
		reply_ls[ss] = reply
	}
	
	for e:=id_ls.Front();e!=nil;e=e.Next() {
		ss := e.(string)
		call, _ := wait_ls[ss]
		reply_call := <- call.Done
		reply, _ := reply_ls[ss]
		if reply.Status != airlineproto.OK {
			final_status = reply.Status
			isSuccess = false
		}
	}
	ori_reply.Status = final_status
	if shouldAbort {
		ori_reply.Status = coordproto.ENOFLIGHT
	}
	return nil	
}