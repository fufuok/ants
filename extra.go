package ants

import (
	"sync/atomic"
)

// SetDefaultAntsPool initialize to the default pool.
func SetDefaultAntsPool(size int, options ...Option) {
	defaultAntsPool, _ = NewPool(size, options...)
}

// MaxBlockingTasks returns the maximum number of goroutines that are blocked when it reaches the capacity of default pool.
func MaxBlockingTasks() int {
	return defaultAntsPool.MaxBlockingTasks()
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked when it reaches the capacity of default pool.
func TuneMaxBlockingTasks(size int) {
	defaultAntsPool.TuneMaxBlockingTasks(size)
}

// MaxBlockingTasks returns the maximum number of goroutines that are blocked when it reaches the capacity of pool.
func (p *poolCommon) MaxBlockingTasks() int {
	return int(atomic.LoadInt32(&p.options.MaxBlockingTasks))
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked when it reaches the capacity of pool.
func (p *poolCommon) TuneMaxBlockingTasks(size int) {
	atomic.StoreInt32(&p.options.MaxBlockingTasks, int32(size))
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPool) TuneMaxBlockingTasks(size int) {
	for _, pool := range mp.pools {
		pool.TuneMaxBlockingTasks(size)
	}
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPoolWithFunc) TuneMaxBlockingTasks(size int) {
	for _, pool := range mp.pools {
		pool.TuneMaxBlockingTasks(size)
	}
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPoolWithFuncGeneric[T]) TuneMaxBlockingTasks(size int) {
	for _, pool := range mp.pools {
		pool.TuneMaxBlockingTasks(size)
	}
}
