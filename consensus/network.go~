// LSP network handler
// Handles Sending and Receiving UDP messages. 
// Marshal and Unmarshal between LspMessage and []byte are also performed here
// by Yu Su <ysu1@andrew.cmu.edu>

package consensus

import (
    "encoding/json"
    "strconv"
)

// Received Data
type NetRead struct {
    msg  *LspMessage
    addr *lspnet.UDPAddr
}

// Network Handler
type iNetworkHandler struct {
    port       int
    conn       *lspnet.UDPConn
    addr       *lspnet.UDPAddr
    ReadC      *chan *NetRead
    TerminateC chan interface{}
    ClosedC    chan interface{}
}

// Close UDP connection and interrupt the network handler
func (h *iNetworkHandler) Close() {
    h.conn.Close()
}

// Dial to server
func (h *iNetworkHandler) Dial(service string) error {
    addr, err := lspnet.ResolveUDPAddr("udp4", service)
    if err != nil {
        return err
    }
    h.addr = addr
    h.conn, err = lspnet.DialUDP("udp4", nil, addr)
    if err != nil {
        return err
    }
    return nil
}

// Listen on port
func (h *iNetworkHandler) Listen(port int) error {
    server := "127.0.0.1:" + strconv.FormatInt(int64(port), 10)
    addr, err := lspnet.ResolveUDPAddr("udp4", server)
    if err != nil {
        return err
    }
    h.addr = addr
    h.conn, err = lspnet.ListenUDP("udp", h.addr)
    if err != nil {
        return err
    }
    return nil
}

// Marshal LspMessage and send it to server
func (n *iNetworkHandler) Send(msg *LspMessage) error {
    buf, err := json.Marshal(*msg)
    if err != nil {
        return err
    }
    _, err = n.conn.Write(buf)
    if err != nil {
        return err
    }
    return nil
}

// Marshal LspMessage and send it to client
func (n *iNetworkHandler) SendMsg(msg *LspMessage, addr *lspnet.UDPAddr) error {
    buf, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    _, err = n.conn.WriteToUDP(buf, addr)
    if err != nil {
        return err
    }
    return nil
}

// Run the network handler, push UDP message into read channel
func (h *iNetworkHandler) run() {
    buf := make([]byte, 2000)
    for {
        n, addr, err := h.conn.ReadFromUDP(buf[0:])
        select {
        case <-h.TerminateC:
            h.conn.Close()
            h.ClosedC <- nil
            return
        default:
            if err == nil {
                data := &NetRead{}
                data.msg = &LspMessage{}
                err := json.Unmarshal(buf[:n], data.msg)
                if err != nil {
                    lsplog.Vlogf(0, "Unmarshal err: %s\n", err.Error())
                } else {
                    data.addr = addr
                    *h.ReadC <- data
                }
            }
        }
    }
}
