// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package common

import (
	"os"

	"github.com/djmitche/dd-agent-comp-experiments/comp/autodiscovery/scheduler"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// SharedOptions defines fx.Options that are shared among all agent flavors.
//
// The confFilePath is passed to the comp/core/config component.
//
// If oneShot is true, then this is a "one-shot" process and all support for long-term
// execution, such as health monitoring, will be disabled.
func SharedOptions(confFilePath string, oneShot bool) fx.Option {
	options := []fx.Option{}

	options = append(options,
		fx.Supply(&core.BundleParams{
			AutoStart:    !oneShot,
			ConfFilePath: confFilePath,
			Console:      true,
		}),
		core.Bundle)

	options = append(options,
		scheduler.Module)

	// Include Fx's detailed logging to stderr only if TRACE_FX is set.
	// This logging is verbose, and occurs mostly during early application
	// startup, before the log component is ready to handle logs.
	options = append(options,
		fx.WithLogger(
			func() fxevent.Logger {
				if os.Getenv("TRACE_FX") == "" {
					return fxevent.NopLogger
				} else {
					return &fxevent.ConsoleLogger{W: os.Stderr}
				}
			},
		),
	)

	return fx.Options(options...)
}
