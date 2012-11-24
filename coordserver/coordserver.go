package main

import (
    "../coordimpl"
    "../coordrpc"
    "flag"
    "net"
    "fmt"
    "log"
    "net/http"
    "net/rpc"
    "strconv"
)

var portnum *int = flag.Int("port", 0, "port # to listen on. nodes default to using an ephemeral port (0).")
var path *string = flag.String("path","config/config","the path of the config file.")

func main() {
    flag.Parse()
    l, e := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
    if e != nil {
        log.Fatal("listen error:", e)
    }
    _, listenport, _ := net.SplitHostPort(l.Addr().String())
    log.Println("Server starting on ", listenport)
    *portnum, _ = strconv.Atoi(listenport)
    co := coordimpl.NewCoordinator(*path)

    corpc := coordrpc.NewCoordinatorRPC(co)
    rpc.Register(corpc)
    rpc.HandleHTTP()
    http.Serve(l, nil)
}
