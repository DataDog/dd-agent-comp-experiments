// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package main

import (
	"go.uber.org/fx"

	"github.com/djmitche/dd-agent-comp-experiments/cmd"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/file"
)

func logsAgentPluginOptions() fx.Option {
	return fx.Options(
		// this list would be different for other agent flavors
		file.Module,
		fx.Invoke(func(file.Component) {}),
	)
}

func main() {
	app := fx.New(
		cmd.SharedOptions("/etc/datadog-agent/datadog.yaml"),
		logs.Module,
		logsAgentPluginOptions(),
	)
	app.Run()
}
