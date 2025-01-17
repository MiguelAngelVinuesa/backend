package utils

import (
	"context"
	"sync"
)

// Finalizer represents the signalling interface for stopping backgrounbd tasks.
type Finalizer interface {
	Done() context.Context
	Shutdown()
}

// Final returns the global Finalizer.
func Final() Finalizer {
	once.Do(func() {
		final = &finalizer{}
		final.done, final.shutdown = context.WithCancel(context.Background())
	})
	return final
}

var (
	final *finalizer
	once  sync.Once
)

type finalizer struct {
	done     context.Context
	shutdown func()
}

// Done returns the finalizer context for background tasks to attach to.
func (f *finalizer) Done() context.Context {
	return f.done
}

// Shutdown cancels the finalizer context to signal attached background tasks to stop.
func (f *finalizer) Shutdown() {
	f.shutdown()
}
