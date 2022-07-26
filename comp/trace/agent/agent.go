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

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/status"
	"github.com/DataDog/dd-agent-comp-experiments/comp/trace/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/trace/internal/httpreceiver"
	"go.uber.org/fx"
)

type agent struct {
	log log.Component
}

type dependencies struct {
	fx.In

	Lc           fx.Lifecycle
	Params       internal.BundleParams
	Config       config.Component
	HTTPReceiver httpreceiver.Component // required just to load the component
	Log          log.Component
}

func newAgent(deps dependencies) (Component, status.Registration) {
	// TODO: this will likely carry a reference to Receiver, Processor, and so
	// on to handle requests for Status, stats, etc.
	a := &agent{
		log: deps.Log,
	}

	var reg status.Registration
	if deps.Params.ShouldStart(deps.Config) {
		deps.Lc.Append(fx.Hook{OnStart: a.start, OnStop: a.stop})
		reg = status.NewRegistration("trace-agent", 3, a.status)
	}

	return a, reg
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
