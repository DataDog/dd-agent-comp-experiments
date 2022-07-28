// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package httpreceiver listens for incoming spans via HTTP and submits them to
// the APM agent pipeline.
package httpreceiver

import (
	"go.uber.org/fx"
)

// team: trace-agent

// Component is the component type.
type Component interface {
	// Enable enables startup of this component.  If not enabled, the component
	// will not start.
	Enable()
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	"comp/trace/internal/httpreceiver",
	fx.Provide(newReceiver),
)
