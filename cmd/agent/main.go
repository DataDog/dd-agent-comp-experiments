// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package main

import (
	"os"

	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/flare"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/health"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/root"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/run"
	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/status"
)

func main() {
	cmd := root.MakeCommand(
		run.Cmd,
		health.Cmd,
		flare.Cmd,
		status.Cmd,
	)
	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
