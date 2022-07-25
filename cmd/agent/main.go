package main

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/agent"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/file"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/manager"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
)

func main() {
	app := fx.New(
		// shared
		log.Module,
		config.Module("/etc/datadog-agent/datadog.yaml"), // XXX this would come from a CLI arg

		// logs-agent
		agent.Module,
		manager.Module,
		file.Module,

		// Invoke just has to require the top-level components and plug-ins
		fx.Invoke(
			func(
				// top-level components
				agent.Component,

				// plugins
				file.Component,
			) {
			}),

		// TODO: This will probably be global to all binaries, so maybe cmd.Module?
		fx.WithLogger(
			func() fxevent.Logger {
				// (we'd probably want to hook this into agent logging at trace level)
				return &fxevent.ConsoleLogger{W: os.Stderr}
			},
		),
	)
	app.Run()
}
