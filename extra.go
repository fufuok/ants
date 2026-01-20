package ants

import (
	"sync/atomic"
)

// SwapDefaultAntsPool sets the default pool to the given pool instance and returns the old pool.
// This function allows you to replace the default pool with any implementation of Pooler,
// such as a MultiPool, MultiPoolWithFunc, or your own custom pool implementation.
//
// If the provided pool is nil, the default pool remains unchanged and nil is returned
// to prevent users from accessing the system default pool for unexpected operations.
//
// The old default pool is returned to the caller only when a valid pool is provided,
// and the caller is responsible for managing its lifecycle.
// The caller should release the old pool when it's no longer needed to prevent goroutine leaks.
//
// Example:
//
//	// Create a MultiPool and set it as the default pool
//	mp, _ := ants.NewMultiPool(2, 1000, ants.RoundRobin)
//	oldPool := ants.SwapDefaultAntsPool(mp)
//
//	// Release the old pool when it's no longer needed
//	if oldPool != nil {
//		oldPool.Release()
//	}
//
// Note: This function is not thread-safe. It should be called during application initialization,
// before any goroutines are started or any tasks are submitted to the default pool.
func SwapDefaultAntsPool(pool Pooler) Pooler {
	// If the provided pool is nil, don't return the old default pool to prevent unauthorized access
	if pool == nil {
		return nil
	}
	oldPool := defaultAntsPool
	defaultAntsPool = pool
	return oldPool
}

// SetDefaultPool safely sets the default pool to the given pool instance and releases the old pool.
// This function is a safer alternative to SwapDefaultAntsPool for most use cases,
// as it automatically releases the old pool to prevent goroutine leaks.
//
// If the provided pool is nil, the default pool remains unchanged and no action is taken.
//
// The function ensures that the old pool's resources are properly released,
// including stopping its cleanup goroutines, heartbeat goroutines, and worker queue.
//
// Example:
//
//	// Create a MultiPool and set it as the default pool
//	mp, _ := ants.NewMultiPool(2, 1000, ants.RoundRobin)
//	ants.SetDefaultPool(mp)
//
//	// The old pool is automatically released, no need to manually release it
//
// Note: This function is not thread-safe. It should be called during application initialization,
// before any goroutines are started or any tasks are submitted to the default pool.
func SetDefaultPool(pool Pooler) {
	if oldPool := SwapDefaultAntsPool(pool); oldPool != nil && oldPool != pool {
		oldPool.Release()
	}
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

// IdleWorkers returns the number of workers currently idle in the pool.
// This includes workers that are available to handle new tasks but are not currently running any tasks.
//
// The idle workers are stored in the workerQueue and are ready to be reused.
// This value plus the Running() value gives the total number of workers in the pool.
func (p *poolCommon) IdleWorkers() int {
	p.lock.Lock()
	idle := p.workers.len()
	p.lock.Unlock()
	return idle
}

// TotalWorkers returns the total number of workers in the pool, including both running and idle workers.
// This is the sum of Running() and IdleWorkers().
//
// This value represents the total number of goroutines that have been created by the pool
// and are either running tasks or waiting to be reused for new tasks.
func (p *poolCommon) TotalWorkers() int {
	return p.Running() + p.IdleWorkers()
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

// IdleWorkers returns the number of workers currently idle in all pools of the multi-pool.
// This includes workers that are available to handle new tasks but are not currently running any tasks.
//
// The idle workers are stored in the workerQueue of each sub-pool and are ready to be reused.
// This value plus the Running() value gives the total number of workers in the multi-pool.
func (mp *MultiPool) IdleWorkers() int {
	idle := 0
	for _, pool := range mp.pools {
		idle += pool.IdleWorkers()
	}
	return idle
}

// TotalWorkers returns the total number of workers in all pools of the multi-pool,
// including both running and idle workers.
// This is the sum of Running() and IdleWorkers().
//
// This value represents the total number of goroutines that have been created by all sub-pools
// and are either running tasks or waiting to be reused for new tasks.
func (mp *MultiPool) TotalWorkers() int {
	return mp.Running() + mp.IdleWorkers()
}

// IdleWorkers returns the number of workers currently idle in all pools of the multi-pool.
// This includes workers that are available to handle new tasks but are not currently running any tasks.
//
// The idle workers are stored in the workerQueue of each sub-pool and are ready to be reused.
// This value plus the Running() value gives the total number of workers in the multi-pool.
func (mp *MultiPoolWithFunc) IdleWorkers() int {
	idle := 0
	for _, pool := range mp.pools {
		idle += pool.IdleWorkers()
	}
	return idle
}

// TotalWorkers returns the total number of workers in all pools of the multi-pool,
// including both running and idle workers.
// This is the sum of Running() and IdleWorkers().
//
// This value represents the total number of goroutines that have been created by all sub-pools
// and are either running tasks or waiting to be reused for new tasks.
func (mp *MultiPoolWithFunc) TotalWorkers() int {
	return mp.Running() + mp.IdleWorkers()
}

// IdleWorkers returns the number of workers currently idle in all pools of the multi-pool.
// This includes workers that are available to handle new tasks but are not currently running any tasks.
//
// The idle workers are stored in the workerQueue of each sub-pool and are ready to be reused.
// This value plus the Running() value gives the total number of workers in the multi-pool.
func (mp *MultiPoolWithFuncGeneric[T]) IdleWorkers() int {
	idle := 0
	for _, pool := range mp.pools {
		idle += pool.IdleWorkers()
	}
	return idle
}

// TotalWorkers returns the total number of workers in all pools of the multi-pool,
// including both running and idle workers.
// This is the sum of Running() and IdleWorkers().
//
// This value represents the total number of goroutines that have been created by all sub-pools
// and are either running tasks or waiting to be reused for new tasks.
func (mp *MultiPoolWithFuncGeneric[T]) TotalWorkers() int {
	return mp.Running() + mp.IdleWorkers()
}

// TuneMaxBlockingTasks changes the maximum number of goroutines that are blocked each pool in multi-pool.
func (mp *MultiPoolWithFuncGeneric[T]) TuneMaxBlockingTasks(size int) {
	perSize := size / len(mp.pools)
	for _, pool := range mp.pools {
		pool.TuneMaxBlockingTasks(perSize)
	}
}
