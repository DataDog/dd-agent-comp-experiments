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
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
)

type agent struct {
	log log.Component
}

type dependencies struct {
	fx.In

	Lc           fx.Lifecycle
	HTTPReceiver httpreceiver.Component // required just to load the component
	Log          log.Component
}

type provides struct {
	fx.Out

	Component
	StatusReg *status.Registration `group:"status"`
}

func newAgent(deps dependencies) provides {
	// TODO: this will likely carry a reference to Receiver, Processor, and so
	// on to handle requests for Status, stats, etc.
	a := &agent{
		log: deps.Log,
	}

	deps.Lc.Append(fx.Hook{
		OnStart: a.start,
		OnStop:  a.stop,
	})

	return provides{
		Component: a,
		StatusReg: status.NewRegistration("trace-agent", 3, a.status),
	}
}

func (a *agent) start(context.Context) error {
	a.log.Debug("Starting trace-agent")
	return nil
}

func (a *agent) stop(context.Context) error {
	a.log.Debug("Stopping trace-agent")
	return nil
}

func (a *agent) status() string {
	var bldr strings.Builder

	fmt.Fprintf(&bldr, "===========\n")
	fmt.Fprintf(&bldr, "Trace Agent\n")
	fmt.Fprintf(&bldr, "===========\n")
	fmt.Fprintf(&bldr, "\n")
	fmt.Fprintf(&bldr, "STATUS: Doin' just fine, thanks!\n")

	return bldr.String()
}
