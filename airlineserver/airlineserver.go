package main

import (
    "../airlineimpl"
    "../airlinerpc"
    "flag"
    "net"
    "fmt"
    "log"
    "net/http"
    "net/rpc"
    "strconv"
)

var portnum *int = flag.Int("port", 0, "port # to listen on. nodes default to using an ephemeral port (0).")

func main() {
    flag.Parse()
    l, e := net.Listen("tcp", fmt.Sprintf(":%d", *portnum))
    if e != nil {
        log.Fatal("listen error:", e)
    }
    _, listenport, _ := net.SplitHostPort(l.Addr().String())
    log.Println("Server starting on ", listenport)
    *portnum, _ = strconv.Atoi(listenport)
    as := airlineimpl.NewAirlineServer()

    asrpc := airlinerpc.NewAirlineServerRPC(as)
    rpc.Register(asrpc)
    rpc.HandleHTTP()
    http.Serve(l, nil)
}
