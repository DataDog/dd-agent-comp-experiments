package main

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/agent"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
)

func main() {
	app := fx.New(
		log.Module,
		agent.Module,
		config.Module("/etc/datadog-agent/datadog.yaml"), // XXX this would come from a CLI arg

		// Invoke just has to require the top-level components
		fx.Invoke(func(agent.Component) {}),

		// This will probably be global to all binaries, so maybe cmd.Module?
		fx.WithLogger(
			func() fxevent.Logger {
				// (we'd probably want to hook this into agent logging at trace level)
				return &fxevent.ConsoleLogger{W: os.Stderr}
			},
		),
	)
	app.Run()
}
