// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package agent implements a component representing the trace agent.
package agent

import (
	"go.uber.org/fx"
)

// team: trace-agent

// Component is the component type.
type Component interface {
	// Enable enables startup of this agent.  If not enabled, the agent will
	// not start.
	Enable()
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	"comp/trace/agent",
	fx.Provide(newAgent),
)
