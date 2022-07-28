// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package logs collects the packages related to the logs agent.
package trace

import (
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/agent"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/httpreceiver"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/processor"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/tracewriter"
	"go.uber.org/fx"
)

// Module defines the fx options for the trace agent.
var Module fx.Option = fx.Module(
	"comp/trace",
	agent.Module,
	processor.Module,
	tracewriter.Module,
	httpreceiver.Module,

	// load and start the top-level agent component
	fx.Invoke(func(agent.Component) {}),
)
