// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package log

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/fx"
)

type logger struct {
	console bool
	level   string
}

type dependencies struct {
	fx.In

	Lc     fx.Lifecycle
	Params ModuleParams `optional:"true"`
}

func newLogger(deps dependencies) Component {
	c := &logger{
		console: deps.Params.Console,
	}
	deps.Lc.Append(fx.Hook{OnStart: c.start})
	return c
}

// Configure implements Component#Configure.
func (c *logger) Configure(level string) error {
	if c.level != "" {
		return errors.New("Do not call Configure() twice, nor after startup")
	}
	c.level = level
	return nil
}

func (c *logger) start(context.Context) error {
	// apply defaults if Configure wasn't called
	if c.level == "" {
		c.level = "warn"
	}

	// (set up seelog with the given level)

	return nil
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
