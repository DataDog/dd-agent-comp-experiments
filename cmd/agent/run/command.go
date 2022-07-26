// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package run implements the `agent run` command.
package run

import (
	"github.com/DataDog/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/DataDog/dd-agent-comp-experiments/cmd/common"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/launchers/file"
	"github.com/DataDog/dd-agent-comp-experiments/comp/trace"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/fxapps"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/startup"
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
	return fxapps.Run(
		common.SharedOptions(root.ConfFilePath, false),
		fx.Supply(logs.BundleParams{
			AutoStart: startup.IfConfigured,
		}),
		logs.Bundle,
		fx.Supply(trace.BundleParams{
			AutoStart: startup.IfConfigured,
		}),
		trace.Bundle,
		logsAgentPluginOptions(),
	)
}
