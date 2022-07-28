// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package run implements the `agent run` command.
package run

import (
	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/common"
	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs"
	logsagent "github.com/djmitche/dd-agent-comp-experiments/comp/logs/agent"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/file"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace"
	traceagent "github.com/djmitche/dd-agent-comp-experiments/comp/trace/agent"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	Cmd = &cobra.Command{
		Use:   "run",
		Short: "Run the Agent",
		Long:  `Runs the agent in the foreground`,
		RunE:  run,
	}
)

func logsAgentPluginOptions() fx.Option {
	return fx.Options(
		// this list would be different for other agent flavors
		file.Module,
		fx.Invoke(func(file.Component) {}),
	)
}

func run(_ *cobra.Command, args []string) error {
	app := fx.New(
		common.SharedOptions(root.ConfFilePath, false),
		logs.Module,
		trace.Module,
		logsAgentPluginOptions(),
		fx.Invoke(func(config config.Component, agent logsagent.Component) {
			if config.GetBool("logs_enabled") {
				agent.Enable()
			}
		}),
		fx.Invoke(func(config config.Component, agent traceagent.Component) {
			if config.GetBool("apm_config.enabled") {
				agent.Enable()
			}
		}),
	)
	return common.RunApp(app)
}
