// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package processor handles processing spans for the trace agent.
//
// It operates a fleet of workers to process inputs concurrently.
package processor

import (
	"github.com/djmitche/dd-agent-comp-experiments/pkg/trace/api"
	"go.uber.org/fx"
)

// team: trace-agent

// Component is the component type.
type Component interface {
	// PayloadChan returns the channel to which receiver components should direct Payloads.
	PayloadChan() chan<- *api.Payload
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	"comp/trace/internal/processor",
	fx.Provide(newProcessor),
)
