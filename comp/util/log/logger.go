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
	level string
}

func newLogger(lc fx.Lifecycle) Component {
	c := &logger{}
	lc.Append(fx.Hook{OnStart: c.start})
	return c
}

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

// Debug implements Logger#Debug.
func (*logger) Debug(v ...interface{}) {
	// stand-in, to avoid messing with seelog
	fmt.Println(v...)
}

// Flush implements Logger#Flush.
func (*logger) Flush() {
	// do nothing
}
