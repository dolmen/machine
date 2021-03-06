package machine

import (
	"context"
	"sync"
	"time"
)

// Routine is an interface representing a goroutine
type Routine interface {
	// Context returns the goroutines unique context that may be used for cancellation
	Context() context.Context
	// Cancel cancels the context returned from Context()
	Cancel()
	// PID() is the goroutines unique process id
	PID() int
	// Tags() are the tags associated with the goroutine
	Tags() []string
	// Start is when the goroutine started
	Start() time.Time
	// Duration is the duration since the goroutine started
	Duration() time.Duration
	// Publish publishes the object to the given channel
	Publish(channel string, obj interface{}) error
	// Subscribe subscribes to a channel and executes the function on every message passed to it. It exits if the goroutines context is cancelled.
	Subscribe(channel string, handler func(obj interface{})) error
	// Machine returns the underlying routine's machine instance
	Machine() *Machine
}

type goRoutine struct {
	machine  *Machine
	ctx      context.Context
	id       int
	tags     []string
	start    time.Time
	doneOnce sync.Once
	cancel   func()
}

func (r *goRoutine) Context() context.Context {
	return r.ctx
}

func (r *goRoutine) PID() int {
	return r.id
}

func (r *goRoutine) Tags() []string {
	return r.tags
}

func (r *goRoutine) Cancel() {
	r.cancel()
}

func (r *goRoutine) Start() time.Time {
	return r.start
}

func (r *goRoutine) Duration() time.Duration {
	return time.Since(r.start)
}

func (g *goRoutine) Publish(channel string, obj interface{}) error {
	return g.machine.pubsub.Publish(channel, obj)
}

func (g *goRoutine) Subscribe(channel string, handler func(obj interface{})) error {
	return g.machine.pubsub.Subscribe(g.ctx, channel, handler)
}

func (g *goRoutine) Machine() *Machine {
	return g.machine
}

func (g *goRoutine) done() {
	g.doneOnce.Do(func() {
		g.cancel()
		g.machine.mu.Lock()
		delete(g.machine.routines, g.id)
		g.machine.mu.Unlock()
	})
}
