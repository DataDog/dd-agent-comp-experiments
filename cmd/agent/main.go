// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package main

import (
	"os"

	"github.com/djmitche/dd-agent-comp-experiments/cmd/agent/root"
	_ "github.com/djmitche/dd-agent-comp-experiments/cmd/agent/run"
)

func main() {
	if err := root.AgentCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
