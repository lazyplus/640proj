package consensus

import (
    "../paxosproto"
    "net"
    "encoding/json"
    "fmt"
)

// Network Handler
type iNetworkHandler struct {
    port       int
    conn       *net.UDPConn
    addr       *net.UDPAddr
    ReadC      *chan *paxosproto.Packet
}

// Close UDP connection and interrupt the network handler
func (h *iNetworkHandler) Close() {
    h.conn.Close()
}

// Listen on port
func (h *iNetworkHandler) Listen(hostport string) error {
    server := hostport
    addr, err := net.ResolveUDPAddr("udp4", server)
    if err != nil {
        fmt.Println(err)
        return err
    }
    h.addr = addr
    h.conn, err = net.ListenUDP("udp4", h.addr)
    if err != nil {
        fmt.Println(err)
        return err
    }
    return nil
}

// Marshal LspMessage and send it to client
func (n *iNetworkHandler) SendMsg(msg *paxosproto.Packet, addr *net.UDPAddr) error {
    buf, err := json.Marshal(msg)
    if err != nil {
        fmt.Println(err)
        return err
    }
    _, err = n.conn.WriteToUDP(buf, addr)
    if err != nil {
        fmt.Println(err)
        return err
    }
    return nil
}

// Run the network handler, push UDP message into read channel
func (h *iNetworkHandler) run() {
    buf := make([]byte, 2000)
    for {
        n, _, err := h.conn.ReadFromUDP(buf[0:])
        if err == nil {
            data := &paxosproto.Packet{}
            err := json.Unmarshal(buf[:n], data)
            if err != nil {
            } else {
                *h.ReadC <- data
            }
        }
    }
}
