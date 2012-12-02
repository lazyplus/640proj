package consensus

type Delegate struct {
    Action interface{}
    timer chan interface{}
}

func NewDelegate(path string) *Delegate {
    dg := &Delegate{}
    return dg
}

func (dg *Delegate) PrepareCancelFlight(args *airlineproto.BookArgs, reply *airlineproto.BookReply) error {
    // make up Value

    *reply = *dg.Push(V)
    return nil
}

func (dg *Delegate) Push(V ValueStruct) interface{} {
    index := 0

    for ; ; index = (index + 1) % len(dg.servers) {
        if dg.prepareRPC(index) {
            continue
        }
        var reply interface{}
        dg.cli[index].Call("AirlineServerRPC.Propose", V, reply)
        if reply != nil {
            return reply
        }
        timer.Sleep(time.Second)
    }
}
