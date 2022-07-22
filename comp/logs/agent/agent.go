// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package agent

import (
	"context"

	"go.uber.org/fx"

	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
)

type agent struct {
	cfg *config
	log log.Component
}

func newAgent(lc fx.Lifecycle, cfg *config, log log.Component) Component {
	c := &agent{
		cfg: cfg,
		log: log,
	}
	lc.Append(fx.Hook{
		OnStart: c.start,
		OnStop:  c.stop,
	})
	return c
}

func (c *agent) start(context.Context) error {
	c.log.Debug("Starting logs-agent")
	c.log.Debug("CCA setting:", c.cfg.containerCollectAll)
	return nil
}

func (c *agent) stop(context.Context) error {
	c.log.Debug("Stopping logs-agent")
	return nil
}
