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
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/agent"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	Cmd = &cobra.Command{
		Use:   "run",
		Short: "Run the Trace-Agent",
		Long:  `Runs the trace-agent in the foreground`,
		RunE:  run,
	}
)

func run(_ *cobra.Command, args []string) error {
	app := fx.New(
		common.SharedOptions(root.ConfFilePath, false),
		trace.Module,
		fx.Invoke(func(config config.Component, agent agent.Component) {
			// enable the trace-agent unconditionally in this binary,
			// regardless of apm_config.enabled (XXX we don't have to do this
			// -- just demonstrating that it's possible)
			agent.Enable()
		}),
	)
	return common.RunApp(app)
}
