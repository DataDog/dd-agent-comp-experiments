// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package agent implements a component representing the trace agent.
package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/djmitche/dd-agent-comp-experiments/comp/status"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/httpreceiver"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/processor"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
)

type agent struct {
	enabled      bool
	httpReceiver httpreceiver.Component
	processor    processor.Component
	log          log.Component
}

type dependencies struct {
	fx.In

	Lc           fx.Lifecycle
	HTTPReceiver httpreceiver.Component
	Processor    processor.Component
	Status       status.Component
	Log          log.Component
}

func newAgent(deps dependencies) Component {
	a := &agent{
		httpReceiver: deps.HTTPReceiver,
		processor:    deps.Processor,
		log:          deps.Log,
	}

	deps.Status.RegisterSection("trace-agent", 3, a.status)

	deps.Lc.Append(fx.Hook{
		OnStart: a.start,
		OnStop:  a.stop,
	})

	return a
}

// Enable implements Component#Enable.
func (a *agent) Enable() {
	a.enabled = true
	a.httpReceiver.Enable()
	a.processor.Enable()
}

func (a *agent) start(context.Context) error {
	if a.enabled {
		a.log.Debug("Starting trace-agent")
	}
	return nil
}

func (a *agent) stop(context.Context) error {
	if a.enabled {
		a.log.Debug("Stopping trace-agent")
	}
	return nil
}

func (a *agent) status() string {
	var bldr strings.Builder

	fmt.Fprintf(&bldr, "===========\n")
	fmt.Fprintf(&bldr, "Trace Agent\n")
	fmt.Fprintf(&bldr, "===========\n")
	fmt.Fprintf(&bldr, "\n")
	if !a.enabled {
		fmt.Fprintf(&bldr, "disabled\n")
		return bldr.String()
	}

	fmt.Fprintf(&bldr, "STATUS: Doin' just fine, thanks!\n")

	return bldr.String()
}
