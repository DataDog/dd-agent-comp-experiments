// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package run implements the `agent run` command.
package run

import (
	"fmt"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/common"
	"github.com/djmitche/dd-agent-comp-experiments/comp/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/file"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the Agent",
		Long:  `Runs the agent in the foreground`,
		RunE:  run,
	}
)

func init() {
	root.AgentCmd.AddCommand(runCmd)
}

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
		logsAgentPluginOptions(),
		// XXX temporary
		fx.Invoke(func(health health.Component) {
			go func() {
				time.Sleep(time.Second / 2)
				fmt.Printf("health:%#v\n", health.GetHealth())
			}()
		}),
	)
	app.Run()
	return app.Err()
}
