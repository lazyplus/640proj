package ultility

import(
    "sync"
)

type PLocke struct {
    cnt_lock sync.Mutex
    cnt int
}

func (l *PLocke) TryLock() bool {
    l.cnt_lock.Lock()
    defer l.cnt_lock.Unlock()

    if l.cnt == 0 {
        l.cnt ++
        return true
    }
    return false
}

func (l* PLocke) SpinLock() {
    for ret := l.TryLock(); !ret; ret = l.TryLock() {
    }
    return
}

func (l *PLocke) Unlock() {
    l.cnt_lock.Lock()
    defer l.cnt_lock.Unlock()

    -- l.cnt
    return
}
