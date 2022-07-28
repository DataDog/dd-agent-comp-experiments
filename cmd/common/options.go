// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package common

import (
	"os"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/flare"
	"github.com/djmitche/dd-agent-comp-experiments/comp/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/ipcapi"
	"github.com/djmitche/dd-agent-comp-experiments/comp/status"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// SharedOptions defines fx.Options that are shared among all agent flavors.
//
// The confFilePath is passed to the comp/config component.
//
// If oneShot is true, then this is a "one-shot" process and all support for long-term
// execution, such as health monitoring, will be disabled.
func SharedOptions(confFilePath string, oneShot bool) fx.Option {
	return fx.Options(
		fx.Supply(log.ModuleParams{Console: !oneShot}),
		log.Module,

		fx.Supply(config.ModuleParams{ConfFilePath: confFilePath}),
		config.Module,

		fx.Supply(health.ModuleParams{Disabled: oneShot}),
		health.Module,

		fx.Supply(ipcapi.ModuleParams{Disabled: oneShot}),
		ipcapi.Module,

		flare.Module,
		status.Module,

		// Include Fx's detailed logging to stderr only if TRACE_FX is set.
		// This logging is verbose, and occurs mostly during early application
		// startup, before the log component is ready to handle logs.
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
}
