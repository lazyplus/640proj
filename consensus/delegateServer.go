package main

import (
	"consensus"
    "flag"
    "net"
    "fmt"
    "log"
    "net/http"
    "net/rpc"
    "strconv"
)

var portnum *int = flag.Int("port", 12340, "port # to listen on. nodes default to using an ephemeral port (0).")
var airline_name string = flag.Int("name","","The airline name that the delegate server belonging to.")
var path string = flag.Int("path","","The path that config files are.")

func main() {
    flag.Parse()
    l, e := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
    if e != nil {
        log.Fatal("listen error:", e)
    }
    _, listenport, _ := net.SplitHostPort(l.Addr().String())
    log.Println("Server starting on ", listenport)
    *portnum, _ = strconv.Atoi(listenport)
    ds := consensus.NewDelegate(path, airline_name, port)

    dsrpc := consensus.NewDelegateServerRPC(ds)
    rpc.Register(dsrpc)
    rpc.HandleHTTP()
    http.Serve(l, nil)
}
