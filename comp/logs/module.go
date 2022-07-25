// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package logs collects the packages related to the logs agent.
package logs

import (
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/agent"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal/sourcemgr"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/launchermgr"
	"go.uber.org/fx"
)

// Module defines the fx options for the logs agent.
//
// It includes all "core" logs-agent components, but does not include launchers.  Applications
// must include the desired launchers itself.
var Module fx.Option = fx.Module(
	"comp/logs",
	agent.Module,
	launchermgr.Module,
	sourcemgr.Module,

	// load and start the top-level agent component
	fx.Invoke(func(agent.Component) {}),
)
