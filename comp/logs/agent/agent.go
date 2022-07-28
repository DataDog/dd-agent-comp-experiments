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

	"github.com/djmitche/dd-agent-comp-experiments/comp/status"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
)

type agent struct {
	enabled bool
	cfg     *config
	log     log.Component
}

type dependencies struct {
	fx.In

	Lc     fx.Lifecycle
	Cfg    *config
	Log    log.Component
	Status status.Component
}

func newAgent(deps dependencies) Component {
	a := &agent{
		cfg: deps.Cfg,
		log: deps.Log,
	}

	deps.Status.RegisterSection("logs-agent", 4, a.status)

	deps.Lc.Append(fx.Hook{
		OnStart: a.start,
		OnStop:  a.stop,
	})

	return a
}

// Enable implements Component#Enable.
func (a *agent) Enable() {
	a.enabled = true
	// TODO: enable subcomponents of this agent
}

func (a *agent) start(context.Context) error {
	if a.enabled {
		a.log.Debug("Starting logs-agent")
	}
	return nil
}

func (a *agent) stop(context.Context) error {
	if a.enabled {
		a.log.Debug("Stopping logs-agent")
	}
	return nil
}

func (a *agent) status() string {
	var bldr strings.Builder

	fmt.Fprintf(&bldr, "==========\n")
	fmt.Fprintf(&bldr, "Logs Agent\n")
	fmt.Fprintf(&bldr, "==========\n")
	fmt.Fprintf(&bldr, "\n")
	if !a.enabled {
		fmt.Fprintf(&bldr, "disabled\n")
		return bldr.String()
	}

	fmt.Fprintf(&bldr, "STATUS: A-OK!\n")

	return bldr.String()
}
