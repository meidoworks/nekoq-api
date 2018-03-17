package workpool

import (
	"log"
	"sync"
	"sync/atomic"

	"goimport.moetang.info/nekoq-api/errorutil"
)

var lock *sync.RWMutex = new(sync.RWMutex)
var poolMap = make(map[string]WorkPool)

var (
	errWorkPoolIsFull = errorutil.New("work pool is full <- workpool <- nekoq-api")
)

type WorkPool interface {
	GetGoroutineNum() int64
	GetMaxGoroutineNum() int64
	Run(string, func()) error
}

type workPoolImp struct {
	maxCnt int64
	curCnt int64
}

func NewOrGetWorkPool(name string, maxPoolSize int64) WorkPool {
	var wp WorkPool
	var ok bool
	lock.RLock()
	wp, ok = poolMap[name]
	lock.RUnlock()
	if ok {
		return wp
	}
	lock.Lock()
	wp, ok = poolMap[name]
	if !ok {
		wp = &workPoolImp{
			maxCnt: maxPoolSize,
			curCnt: 0,
		}
		poolMap[name] = wp
	}
	lock.Unlock()
	return wp
}

func NewUnmanagedWorkPool(name string, maxPoolSize int64) WorkPool {
	log.Println("allocate unmanaged work pool.", "name:", name, "pool size:", maxPoolSize)
	return &workPoolImp{
		maxCnt: maxPoolSize,
		curCnt: 0,
	}
}

func (this *workPoolImp) GetGoroutineNum() int64 {
	return this.curCnt
}

func (this *workPoolImp) GetMaxGoroutineNum() int64 {
	return this.maxCnt
}

func (this *workPoolImp) Run(rName string, f func()) error {
	acquired := false
	max := this.maxCnt
	cnt := atomic.LoadInt64(&this.curCnt)
	if cnt >= max {
		return errWorkPoolIsFull
	}
	for !atomic.CompareAndSwapInt64(&this.curCnt, cnt, cnt+1) {
		cnt = atomic.LoadInt64(&this.curCnt)
		if cnt >= max {
			return errWorkPoolIsFull
		}
	}
	acquired = true
	go func() {
		defer func() {
			if acquired {
				cnt := atomic.LoadInt64(&this.curCnt)
				for !atomic.CompareAndSwapInt64(&this.curCnt, cnt, cnt-1) {
					cnt = atomic.LoadInt64(&this.curCnt)
				}
			}
		}()
		f()
	}()
	return nil
}

func (this *workPoolImp) RunUnsafe(rName string, f func()) error {
	max := this.maxCnt
	cnt := atomic.LoadInt64(&this.curCnt)
	if cnt >= max {
		return errWorkPoolIsFull
	}
	for !atomic.CompareAndSwapInt64(&this.curCnt, cnt, cnt+1) {
		cnt = atomic.LoadInt64(&this.curCnt)
		if cnt >= max {
			return errWorkPoolIsFull
		}
	}
	go func() {
		f()
		cnt := atomic.LoadInt64(&this.curCnt)
		for !atomic.CompareAndSwapInt64(&this.curCnt, cnt, cnt-1) {
			cnt = atomic.LoadInt64(&this.curCnt)
		}
	}()
	return nil
}
