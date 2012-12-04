package ultility

import(
    "sync"
)

type PLock struct {
    cnt_lock sync.Mutex
    cnt int
}

func (l *PLock) TryLock() bool {
    l.cnt_lock.Lock()
    defer l.cnt_lock.Unlock()

    if l.cnt == 0 {
        l.cnt ++
        return true
    }
    return false
}

func (l* PLock) SpinLock() {
    for ret := l.TryLock(); !ret; ret = l.TryLock() {
    }
    return
}

func (l *PLock) Unlock() {
    l.cnt_lock.Lock()
    defer l.cnt_lock.Unlock()

    l.cnt --
    return
}
