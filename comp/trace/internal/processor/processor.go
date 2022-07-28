// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package processor handles processing spans for the trace agent.
package processor

import (
	"context"
	"runtime"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/tracewriter"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/trace/api"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/actor"
	"go.uber.org/fx"
)

// processor implements the singleton controlling the workers.
//
// TODO: actually use workers (perhaps with some general support in pkg/util/chanworkers)
type processor struct {
	// payloadChan is the channel where this component gets the payloads
	// to process
	payloadChan chan *api.Payload

	// traceWriterChan is the channel to which this component writes payloads
	// that should be sent to the Datadog API.
	traceWriterChan chan<- *api.Payload

	// actor implements the actor model for this component
	actor actor.Goroutine

	// health supports monitoring this component
	health *health.ActorRegistration
}

type dependencies struct {
	fx.In

	Lc          fx.Lifecycle
	Health      health.Component
	TraceWriter tracewriter.Component
}

func newProcessor(deps dependencies) Component {
	width := runtime.NumCPU()
	p := &processor{
		payloadChan:     make(chan *api.Payload, width),
		traceWriterChan: deps.TraceWriter.PayloadChan(),
		health:          deps.Health.RegisterActor("comp/trace/internal/processor", 1*time.Second),
	}
	p.actor.HookLifecycle(deps.Lc, p.run)
	return p
}

func (p *processor) PayloadChan() chan<- *api.Payload {
	return p.payloadChan
}

// run implements the component's core loop
func (p *processor) run(ctx context.Context) {
	defer p.health.Stop()
	for {
		select {
		case payload := <-p.payloadChan:
			// XXX there's lots of processing to do here, but it all occurs in
			// the same goroutine via function calls.  That could be done in
			// this component (as it's done in pkg/trace now), or by calling
			// methods on other components.  Using other components would
			// facilitate testing, but otherwise not add a lot of value.

			p.traceWriterChan <- payload
		case <-p.health.Chan():
		case <-ctx.Done():
			return
		}
	}
}
