// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package tracewriter

import (
	"context"
	"time"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/health"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/comp/trace/internal"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/trace/api"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/actor"
	"go.uber.org/fx"
)

type traceWriter struct {
	in chan *api.Payload

	actor  actor.Actor
	log    log.Component
	health *health.Handle
}

type dependencies struct {
	fx.In

	Lc     fx.Lifecycle
	Params internal.BundleParams
	Config config.Component
	Log    log.Component
}

func newTraceWriter(deps dependencies) (Component, health.Registration) {
	healthReg := health.NewRegistration(componentName)
	t := &traceWriter{
		in:     make(chan *api.Payload, 1000),
		log:    deps.Log,
		health: healthReg.Handle,
	}
	if deps.Params.ShouldStart(deps.Config) {
		t.actor.HookLifecycle(deps.Lc, t.run)
	}
	return t, healthReg
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
