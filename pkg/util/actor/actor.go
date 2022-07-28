// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package actor provides basic support for building actors for use in the Agent.
package actor

import (
	"context"

	"go.uber.org/fx"
)

// Goroutine manages an actor goroutine, supporting starting and later stopping
// the goroutine.  This is one-shot: once started and stopped, the goroutine cannot
// be started again.
//
// The zero value is a valid initial state.
type Goroutine struct {
	// started is true after the goroutine has been started, and remains true after
	// it has stopped.
	started bool

	// enabled, if not nil, points to a boolean determining whether this
	// component is enabled.  If not, then start won't start anything.
	enabled *bool

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
func (gr *Goroutine) HookLifecycle(lc fx.Lifecycle, run RunFunc) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			gr.Start(run)
			return nil
		},
		OnStop: gr.Stop,
	})
}

// EnableFlag instructs the goroutine to start only if the pointed-to boolean
// is true.  If this is not called, then the actor starts unconditionally.
func (gr *Goroutine) EnableFlag(enabled *bool) {
	gr.enabled = enabled
}

// Start starts run in a goroutine, setting up to stop it by cancelling the context
// it receives.
func (gr *Goroutine) Start(run RunFunc) {
	if gr.enabled != nil && !*gr.enabled {
		return
	}

	if gr.started {
		panic("Goroutine has already been started")
	}
	gr.started = true

	ctx, cancel := context.WithCancel(context.Background())
	gr.cancel = cancel
	gr.stopped = make(chan struct{})
	go gr.run(run, ctx)
}

// Stop stops the goroutine, waiting until it is complete, or the given context
// is cancelled, before returning.  Returns the error from context if it is
// cancelled.
func (gr *Goroutine) Stop(ctx context.Context) error {
	if gr.enabled != nil && !*gr.enabled {
		return nil
	}

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
func (gr *Goroutine) run(run RunFunc, ctx context.Context) {
	defer close(gr.stopped)
	run(ctx)
}
