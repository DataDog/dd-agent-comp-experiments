// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package tracewriter

import (
	"context"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/trace/api"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/actor"
	"go.uber.org/fx"
)

type traceWriter struct {
	in chan *api.Payload

	actor  actor.Goroutine
	log    log.Component
	health *health.Registration
}

type dependencies struct {
	fx.In

	Lc  fx.Lifecycle
	Log log.Component
}

type provides struct {
	fx.Out

	Component
	HealthReg *health.Registration `group:"health"`
}

func newTraceWriter(deps dependencies) provides {
	t := &traceWriter{
		in:     make(chan *api.Payload, 1000),
		log:    deps.Log,
		health: health.NewRegistration(componentName),
	}
	t.actor.HookLifecycle(deps.Lc, t.run)
	return provides{
		Component: t,
		HealthReg: t.health,
	}
}

func (t *traceWriter) PayloadChan() chan<- *api.Payload {
	return t.in
}

func (t *traceWriter) run(ctx context.Context) {
	monitor, stopMonitor := t.health.LivenessMonitor(time.Second)
	for {
		select {
		case payload := <-t.in:
			t.log.Debug("sending payload", payload)
		case <-monitor:
		case <-ctx.Done():
			stopMonitor()
			return
		}
	}
}
