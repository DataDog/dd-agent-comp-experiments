// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package tracewriter buffers traces and APM events, flushing them to the
// Datadog API.
package tracewriter

import (
	"github.com/DataDog/dd-agent-comp-experiments/pkg/trace/api"
	"go.uber.org/fx"
)

// team: trace-agent

const componentName = "comp/trace/internal/tracewriter"

// Component is the component type.
type Component interface {
	// PayloadChan returns the channel to which components should direct
	// Payloads to be written.
	PayloadChan() chan<- *api.Payload
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	componentName,
	fx.Provide(newTraceWriter),
)
