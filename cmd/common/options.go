// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package common

import (
	"os"

	"github.com/djmitche/dd-agent-comp-experiments/comp/autodiscovery/scheduler"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/flare"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/status"
	"github.com/djmitche/dd-agent-comp-experiments/comp/ipc/ipcclient"
	"github.com/djmitche/dd-agent-comp-experiments/comp/ipc/ipcserver"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
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
		fx.Supply(&log.ModuleParams{Console: !oneShot}),
		log.Module)

	options = append(options,
		fx.Supply(&config.ModuleParams{ConfFilePath: confFilePath}),
		config.Module)

	options = append(options,
		fx.Supply(&health.ModuleParams{Disabled: oneShot}),
		health.Module)

	var ipcInst ipcserver.Component
	options = append(options,
		fx.Supply(&ipcserver.ModuleParams{Disabled: oneShot}),
		fx.Populate(&ipcInst), // instantiate ipc server, even if nothing depends on it
		ipcserver.Module)

	var flareInst flare.Component
	options = append(options,
		fx.Populate(&flareInst), // instantiate flare, even if nothing depends on it
		flare.Module)

	var statusInst status.Component
	options = append(options,
		fx.Populate(&statusInst), // instantiate status, even if nothing depends on it
		status.Module)

	// oneShot processes typically use the ipc client, while 'run' processes do not.
	if oneShot {
		options = append(options,
			ipcclient.Module)
	}

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
