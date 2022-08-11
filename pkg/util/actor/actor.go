// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package actor provides basic support for building actors for use in the Agent.
//
// Methods on this component are not re-entrant.  Components using this one
// should _either_ call HookLifecycle once in their constructor or call Start
// and Stop from their lifecycle hook.
package actor

import (
	"context"

	"go.uber.org/fx"
)

// Actor manages an actor goroutine, supporting starting and later stopping
// the goroutine.  This is one-shot: once started and stopped, the goroutine cannot
// be started again.
//
// The zero value is a valid initial state.
type Actor struct {
	// started is true after the goroutine has been started, and remains true after
	// it has stopped.
	started bool

	// cancel cancels the context passed to the `run` function, used to signal
	// that it should stop
	cancel context.CancelFunc

	// stopped is closed once the run function returns.
	stopped chan struct{}
}

// RunFunc defines the function implementing the actor's event loop.  It should
// run until the passed context is cancelled.
type RunFunc func(context.Context)

// HookLifecycle connects this goroutine to the given fx.Lifecycle, starting and
// stopping it with the lifecycle.  Use this method _or_ the Start and Stop methods,
// but not both.
func (gr *Actor) HookLifecycle(lc fx.Lifecycle, runFunc RunFunc) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			gr.Start(runFunc)
			return nil
		},
		OnStop: gr.Stop,
	})
}

// Start starts run in a goroutine, setting up to stop it by cancelling the context
// it receives.
func (gr *Actor) Start(runFunc RunFunc) {
	if gr.started {
		panic("Goroutine has already been started")
	}
	gr.started = true

	ctx, cancel := context.WithCancel(context.Background())
	gr.cancel = cancel
	gr.stopped = make(chan struct{})
	go gr.run(runFunc, ctx)
}

// Stop stops the goroutine, waiting until it is complete, or the given context
// is cancelled, before returning.  Returns the error from context if it is
// cancelled.
func (gr *Actor) Stop(ctx context.Context) error {
	if !gr.started {
		panic("Goroutine has not been started")
	}
	if gr.cancel == nil {
		panic("Goroutine has already been stopped")
	}
	gr.cancel()
	gr.cancel = nil
	select {
	case <-gr.stopped:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// run executes the given function, ensuring that the stopped channel is closed
// when it finishes.  This method runs in a dedicated goroutine.
func (gr *Actor) run(runFunc RunFunc, ctx context.Context) {
	defer close(gr.stopped)
	runFunc(ctx)
}
