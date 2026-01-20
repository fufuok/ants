package ants_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fufuok/ants"
)

func TestTuneMaxBlockingSubmit(t *testing.T) {
	poolSize := 10
	// Create a new pool with the specified size and options
	pool, err := ants.NewPool(poolSize, ants.WithMaxBlockingTasks(1))
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	// Set the new pool as the default pool
	ants.SetDefaultAntsPool(pool)
	defer func() {
		ants.Release()
		// Create a new default pool and set it as the default
		defaultPool, err := ants.NewPool(ants.DefaultAntsPoolSize)
		if err != nil {
			t.Fatalf("Failed to create default pool: %v", err)
		}
		ants.SetDefaultAntsPool(defaultPool)
	}()
	for i := 0; i < poolSize-1; i++ {
		require.NoError(t, ants.Submit(longRunningFunc), "submit when pool is not full shouldn't return error")
	}
	ch := make(chan struct{})
	f := func() {
		<-ch
	}
	// p is full now.
	require.NoError(t, ants.Submit(f), "submit when pool is not full shouldn't return error")

	// tune max blocking tasks to 3
	ants.TuneMaxBlockingTasks(3)

	var wg sync.WaitGroup
	wg.Add(3)
	errCh := make(chan error, 3)
	for i := 0; i < 3; i++ {
		go func() {
			// should be blocked. blocking num == 3
			if err := ants.Submit(demoFunc); err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}
	time.Sleep(1 * time.Second)
	// already reached max blocking limit
	require.ErrorIsf(t, ants.Submit(demoFunc), ants.ErrPoolOverload,
		"blocking submit when pool reach max blocking submit should return ants.ErrPoolOverload")
	// interrupt f to make blocking submit successful.
	close(ch)
	wg.Wait()
	select {
	case <-errCh:
		t.Fatalf("blocking submit when pool is full should not return error")
	default:
	}
}

func TestTuneMaxBlockingTasksWithFuncGeneric(t *testing.T) {
	poolSize := 10
	p, err := ants.NewPoolWithFuncGeneric(poolSize, longRunningPoolFuncCh, ants.WithMaxBlockingTasks(1))
	require.NoError(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	ch := make(chan struct{})
	for i := 0; i < poolSize-1; i++ {
		require.NoError(t, p.Invoke(ch), "submit when pool is not full shouldn't return error")
	}
	// p is full now.
	require.NoError(t, p.Invoke(ch), "submit when pool is not full shouldn't return error")

	// tune max blocking tasks to 3
	p.TuneMaxBlockingTasks(3)

	var wg sync.WaitGroup
	wg.Add(3)
	errCh := make(chan error, 3)
	for i := 0; i < 3; i++ {
		go func() {
			// should be blocked. blocking num == 3
			if err := p.Invoke(ch); err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}
	time.Sleep(1 * time.Second)
	// already reached max blocking limit
	require.ErrorIsf(t, p.Invoke(ch), ants.ErrPoolOverload,
		"blocking submit when pool reach max blocking submit should return ants.ErrPoolOverload: %v", err)
	// interrupt one func to make blocking submit successful.
	close(ch)
	wg.Wait()
	select {
	case <-errCh:
		t.Fatalf("blocking submit when pool is full should not return error")
	default:
	}
}
