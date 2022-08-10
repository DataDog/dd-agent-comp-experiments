// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package run implements the `agent run` command.
package run

import (
	"github.com/DataDog/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/DataDog/dd-agent-comp-experiments/cmd/common"
	"github.com/DataDog/dd-agent-comp-experiments/comp/trace"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/startup"
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
		fx.Supply(trace.BundleParams{
			AutoStart: startup.IfConfigured,
		}),
		trace.Bundle,
	)
	return common.RunApp(app)
}
