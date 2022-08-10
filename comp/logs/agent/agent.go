// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package agent

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/fx"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/status"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/launchers/launchermgr"
)

type agent struct {
	log         log.Component
	launchermgr launchermgr.Component
}

type dependencies struct {
	fx.In

	Lc          fx.Lifecycle
	Params      internal.BundleParams
	Config      config.Component
	Log         log.Component
	Launchermgr launchermgr.Component
}

type provides struct {
	fx.Out

	Component
	StatusReg *status.Registration `group:"true"`
}

func newAgent(deps dependencies) provides {
	a := &agent{
		log:         deps.Log,
		launchermgr: deps.Launchermgr,
	}

	prov := provides{Component: a}
	if deps.Params.ShouldStart(deps.Config) {
		deps.Lc.Append(fx.Hook{OnStart: a.start, OnStop: a.stop})
		prov.StatusReg = status.NewRegistration("logs-agent", 4, a.status)
	}

	return prov
}

func (a *agent) start(context.Context) error {
	a.log.Debug("Starting logs-agent")
	return nil
}

func (a *agent) stop(context.Context) error {
	a.log.Debug("Stopping logs-agent")
	return nil
}

func (a *agent) status() string {
	var bldr strings.Builder

	fmt.Fprintf(&bldr, "==========\n")
	fmt.Fprintf(&bldr, "Logs Agent\n")
	fmt.Fprintf(&bldr, "==========\n")
	fmt.Fprintf(&bldr, "\n")
	fmt.Fprintf(&bldr, "Running Launchers:\n")

	for name := range a.launchermgr.GetLaunchers() {
		fmt.Fprintf(&bldr, " %s\n", name)
	}

	return bldr.String()
}
