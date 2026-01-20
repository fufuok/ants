package ants

import (
	"sync/atomic"
)

// SetDefaultAntsPool sets the default pool to the given pool instance.
// This function allows you to replace the default pool with any implementation of Pooler,
// such as a MultiPool, MultiPoolWithFunc, or your own custom pool implementation.
//
// Example:
//
//	// Create a MultiPool and set it as the default pool
//	mp, _ := ants.NewMultiPool(2, 1000, ants.RoundRobin)
//	ants.SetDefaultAntsPool(mp)
//
//	// Now all calls to ants.Submit() will use the MultiPool
//
// Note: This function is not thread-safe. It should be called during application initialization,
// before any goroutines are started or any tasks are submitted to the default pool.
func SetDefaultAntsPool(pool Pooler) {
	defaultAntsPool = pool
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

// MaxBlockingTasks returns the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPool) MaxBlockingTasks() int {
	n := 0
	for _, pool := range mp.pools {
		n += pool.MaxBlockingTasks()
	}
	return n
}

// MaxBlockingTasks returns the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPoolWithFunc) MaxBlockingTasks() int {
	n := 0
	for _, pool := range mp.pools {
		n += pool.MaxBlockingTasks()
	}
	return n
}

// MaxBlockingTasks returns the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPoolWithFuncGeneric[T]) MaxBlockingTasks() int {
	n := 0
	for _, pool := range mp.pools {
		n += pool.MaxBlockingTasks()
	}
	return n
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPool) TuneMaxBlockingTasks(size int) {
	perSize := size / len(mp.pools)
	for _, pool := range mp.pools {
		pool.TuneMaxBlockingTasks(perSize)
	}
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPoolWithFunc) TuneMaxBlockingTasks(size int) {
	perSize := size / len(mp.pools)
	for _, pool := range mp.pools {
		pool.TuneMaxBlockingTasks(perSize)
	}
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPoolWithFuncGeneric[T]) TuneMaxBlockingTasks(size int) {
	perSize := size / len(mp.pools)
	for _, pool := range mp.pools {
		pool.TuneMaxBlockingTasks(perSize)
	}
}
