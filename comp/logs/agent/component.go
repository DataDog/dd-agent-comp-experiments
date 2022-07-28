// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package agent implements a component representing the logs agent.  This
// component coordinates activity related to gathering logs and forwarding them
// to the Datadog intake.
package agent

import (
	"go.uber.org/fx"
)

// team: agent-metrics-logs

// Component is the component type.
type Component interface {
	// Enable enables startup of this agent.  If not enabled, the agent will
	// not start.
	Enable()
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	"comp/logs/agent",
	fx.Provide(newConfig),
	fx.Provide(newAgent),
)
