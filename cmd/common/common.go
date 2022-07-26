// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package common

import (
	"os"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// SharedOptions defines fx.Options that are shared among all agent flavors.
func SharedOptions(configFilePath string) fx.Option {
	return fx.Options(
		log.Module,
		config.Module,
		health.Module,
		fx.Invoke(func(cfg config.Component) {
			cfg.Setup(configFilePath)
		}),
		fx.WithLogger(
			func() fxevent.Logger {
				// (we'd probably want to hook this into agent logging at trace level)
				return &fxevent.ConsoleLogger{W: os.Stderr}
			},
		),
	)
}
