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

func loggerOptions() fx.Option {
	return fx.WithLogger(
		func() fxevent.Logger {
			// (we'd probably want to hook this into agent logging at trace level)
			return &fxevent.ConsoleLogger{W: os.Stderr}
		},
	)
}

func sharedOptions(configFilePath string) fx.Option {
	return fx.Options(
		log.Module,
		config.Module,
		fx.Invoke(func(cfg config.Component) {
			cfg.Setup(configFilePath)
		}),
	)
}

func logsAgentOptions() fx.Option {
	return fx.Options(
		agent.Module,
		manager.Module,
		fx.Invoke(func(agent.Component) {}),
	)
}

func logsAgentPluginOptions() fx.Option {
	return fx.Options(
		// this list would be different for other agent flavors
		file.Module,
		fx.Invoke(func(file.Component) {}),
	)
}

func main() {
	app := fx.New(
		loggerOptions(),
		sharedOptions("/etc/datadog-agent/datadog.yaml"),
		logsAgentOptions(),
		logsAgentPluginOptions(),
	)
	app.Run()
}
