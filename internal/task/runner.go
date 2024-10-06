package task

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Runner is used to wait for a group of tasks to finish.
//
// It will stop all the tasks on the first task failure, and the Wait() method will return only the
// first encountered error.
type Runner struct {
	wg             sync.WaitGroup
	ctx            context.Context
	cancelFunc     context.CancelFunc
	firstRunErrPtr unsafe.Pointer
}

// NewTaskRunner creates new task runner instance.
func NewTaskRunner() *Runner {
	ctx, cancel := context.WithCancel(context.Background())

	return &Runner{
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Run runs tasks in the group.
//
// Every task is run in new goroutine.
// When a task returns an error, all the tasks in the group are canceled.
func (r *Runner) Run(tasks ...func(ctx context.Context) error) {
	if r.ctx.Err() != nil {
		return
	}

	for _, fn := range tasks {
		r.wg.Add(1)
		go func(fn func(ctx context.Context) error) {
			defer r.wg.Done()

			err := fn(r.ctx)
			if err != nil {
				r.cancelWithError(err)
			}
		}(fn)
	}
}

// Wait until all tasks are stopped.
func (r *Runner) Wait(ctx context.Context) error {
	if ctx != context.TODO() {
		doneCh := make(chan struct{})
		defer close(doneCh)

		go func() {
			select {
			case <-r.ctx.Done():
			case <-doneCh:
			case <-ctx.Done():
				r.cancelWithError(ctx.Err())
			}
		}()
	}

	r.wg.Wait()

	err := (*error)(atomic.LoadPointer(&r.firstRunErrPtr))
	if err != nil {
		return *err
	}

	return nil
}

func (r *Runner) cancelWithError(err error) {
	swapped := atomic.CompareAndSwapPointer(&r.firstRunErrPtr, nil, (unsafe.Pointer)(&err))

	if swapped {
		r.cancelFunc()
	}
}

// Cancel cancels all the tasks.
func (r *Runner) Cancel() {
	r.cancelFunc()
}
