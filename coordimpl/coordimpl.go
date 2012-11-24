package coordimpl

import {
	"../coordproto"
	"../airlineproto"
	"sync"
	"os"
	"bufio"
	"strings"
	"net/rpc"
	"container/list"
}

type coordserver struct {
	num_airline 	int
	addr_airline 	map[string] string
	map_lock		sync.Mutex
	 
}

func NewCoordinator (path string) *coordserver {
	co := &coordserver{}
	co.addr_airline = make(map[string] string)
	userFile := path
	fin, err := os.Open(userFile)
	defer fin.Close()
	if err != nil{
		return nil
	}
	num := 0
	rf := bufio.NewReader(fin)
	for{
		s, err2 := rf.ReadString('\n')
		if err2 != nil {
			break
		}
		num++
		ss := strings.Split(s,"\t")
		co.addr_airline[ss[0]] = ss[1]
	}
	co.num_airline = num
	return co
}

func (co *coordserver) BookFlights(args *coordproto.BookArgs, ori_reply *coordproto.BookReply) error {
	//2pc
	//assume the flight ID is in format: airline + "::" + ID 
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	//first round
	ls := len(args.Flights)
	wait_ls := List.new()
	reply_ls := List.new()
	client_ls := List.new()
	seq_ls := List.new()
	for i:=0;i<ls;i++ {
		ss := strings.Split(args.Flights[i],"::")
		airline_name = ss[0]
		addr := co.addr_airline[airline_name]
		args_out := &airlineproto.BookArgs{ss[1],args.Email,args.count}
		reply := &airlineproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		bookcall := client.Go("asrpc.PrepareBookFlight",args_out, reply,nil)
		wait_ls.PushBack(bookcall)
		reply_ls.PushBack(reply)
		client_ls.PushBack(client)	
	}
	var shouldCommit int = 1
	var finalstatus int = coordserver.OK
	rt := reply_ls.Front()
	for e := wait_ls.Front();e != nil; e=e.Next() {
		call := e.(*rpc.Call)
		reply_call := <- call.Done
		reply := rt.(*airlineproto.BookReply)
		rt = rt.Next()
		if reply.Status != airlineproto.OK {
			shouldCommit = 0
			finalstatus = reply.Status
		}
		seq_ls.PushBack(reply.Seq)
	}
	
	//second phase
	wait_ls.Init()
	reply_ls.Init()
	se := seq_ls.Front()
	for e:=client_ls.Front();e!=nil;e=e.Next() {
		client := e.(*rpc.Client)
		seq := se.(int)
		se = se.Next()	
		args_out := &airlineproto.DecisionArgs{shouldCommit,seq}
		reply := &airlineproto.DecisionReply{}
		decisioncall := client.Go("asrpc.BookDecision", args_out, reply,nil)
		wait_ls.PushBack(decisioncall)
		reply_ls.PushBack(reply)
	}
	rt = reply.Front()
	
	var isSuccess bool = true
	
	for e:=wait_ls.Front();e!=nil;e=e.Next(){
		call := e.(*rpc.Call)
		reply_call := <- call.Done
		reply := rt.(&airlineproto.DecisionReply)
		rt = rt.Next()
		if reply.Status != airlineproto.OK {
			isSuccess = false
			finalstatus = reply.Status
		}
	}	
	ori_reply.Status = finalstatus
	return nil
}

func (co *coordserver) CancelFlights(args *coordproto.CancelArgs, ori_reply *coordproto.CancelReply) error {
	co.map_lock.Lock()
	defer co.map_lock.Unlock()
	
	ls := len(args.Flights)
	wait_ls := List.new()
	reply_ls := List.new()
	client_ls := List.new()
	seq_ls := List.new()
	//assume the input is airline + "::" + ID
	for i:=0;i<ls;i++ {
		ss := strings.Split(args.Flights[i])
		airline_name := ss[0]
		airline_id := ss[1]
		addr := co.addr_airline[airline_name]
		args_out := &airlineproto.BookArgs{airline_id,args.Email,args.count}
		reply := &airlineproto.BookReply{}
		client, err := rpc.DialHTTP("tcp",addr)
		if err != nil {
			return err
		}
		cancelcall := client.Go("asrpc.repareCancelFlight",args_out,reply,nil)
		client_ls.PushBack(client)
		wait_ls.PushBack(cancelcall)
		reply_ls.PushBack(reply)
	}	
	var should_commit int = 1
	var final_status int = coordproto.OK
	rt := reply.Front()
	for e:=wait_ls.Front();e!=nil;e=e.Next() {
		call := e.(*rpc.Call)
		reply_call := call.Done
		reply := rt.(*airlineproto.BookReply)
		rt = rt.Next()
		if reply.Status != airlineproto.OK {
			should_commit = 0
			final_status = reply.Status
		}
		seq_ls.PushBack(reply.Seq)
	}
	//second phase
	wait_ls.Init()
	reply_ls.Init()
	var isSuccess bool = true
	se := seq_ls.
	for e:=client_ls.Front();e!=nil;e=e.Next {
		client := e.(*rpc.Client)
		seq := se.(int)
		args_out := &airlineproto.DecisionArgs{should_commit,seq}
		reply := &airlineproto.DecisionReply{}
		call := client.Go("asrpc.CancelDecision",args_out,reply)
		wait_ls.PushBack(call)
		reply_ls.PushBack(reply)
	}
	rt := reply.Front()
	
	for e:=wait_ls.Front();e!=nil;e=e.Next() {
		call := e.(*rpc.Call)
		reply_call := <- call.Done
		reply := rt.(*airlineproto.DecisionReply)
		rt = rt.Next()
		if reply.Status != airlineproto.OK {
			final_status = reply.Status
			isSuccess = false
		}
	}
	ori_reply.Status = final_status
	return nil
	
}
