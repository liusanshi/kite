package util

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// SpinLock 自旋锁
type SpinLock uint32

//Lock 锁
func (sl *SpinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		time.Sleep(time.Microsecond * 1)
		runtime.Gosched() //without this it locks up on GOMAXPROCS > 1
	}
}

//Unlock 解锁
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

// TryLock 尝试获取锁
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1)
}

//NewSpinLock 新建一个自旋锁
func NewSpinLock() sync.Locker {
	var lock SpinLock
	return &lock
}
