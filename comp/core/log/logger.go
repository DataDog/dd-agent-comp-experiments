// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package log

import (
	"fmt"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/internal"
	"go.uber.org/fx"
)

type logger struct {
	console bool
	level   string
}

type dependencies struct {
	fx.In

	Lc     fx.Lifecycle
	Params internal.BundleParams
	Config config.Component
}

func newLogger(deps dependencies) Component {
	c := &logger{
		console: deps.Params.Console,
		level:   deps.Config.GetString("log_level"),
	}
	if c.level == "" {
		c.level = "warn"
	}

	return c
}

// Debug implements Component#Debug.
func (l *logger) Debug(v ...interface{}) {
	// stand-in, to avoid messing with seelog
	if l.console {
		fmt.Println(v...)
	}
}

// Flush implements Component#Flush.
func (*logger) Flush() {
	// do nothing
}
