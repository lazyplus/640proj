package main

import (
    "../delegateimpl"
    "../delegateproto"
    "runtime"
    "flag"
    "net"
    "fmt"
    "log"
    "net/http"
    "net/rpc"
    "strconv"
)

var portnum *int = flag.Int("port", 12300, "port # to listen on. nodes default to using an ephemeral port (0).")
var airline_name *string = flag.String("name","1","The airline name that the delegate server belonging to.")
var path *string = flag.String("path","config/config","The path that config files are.")

func main() {
    runtime.GOMAXPROCS(20)
    flag.Parse()
    ds := delegateimpl.NewDelegate(*path, *airline_name)

    l, e := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
    if e != nil {
        log.Fatal("listen error:", e)
    }
    _, listenport, _ := net.SplitHostPort(l.Addr().String())
    log.Println("Server starting on ", listenport)
    *portnum, _ = strconv.Atoi(listenport)

    dsrpc := delegateproto.NewDelegateServerRPC(ds)
    rpc.Register(dsrpc)
    rpc.HandleHTTP()
    http.Serve(l, nil)
}
