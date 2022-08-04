// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package sourcemgr implements a component managing logs-agent sources (type
// LogSource).  It receives additions and removals of sources from other
// components, and it informs subscribers of these additions and removals.
//
// Subscribers should provide a Subscription value in value-group "sourcemgr".
//
// Once added to this component, a LogSource must be considered immutable: neither
// the component having called AddSource, nor any of the subscribers, may modify the
// source.
//
// All component methods can be called concurrently.
package sourcemgr

import (
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/subscriptions"
	"go.uber.org/fx"
)

// team: agent-metrics-logs

const componentName = "comp/logs/internal/sourcemgr"

// Component is the component type.
type Component interface {
	// AddSource adds a new log source.
	AddSource(*LogSource)

	// RemoveSource removes an existing log source.
	RemoveSource(*LogSource)
}

// LogSource defines a source for logs.
type LogSource struct {
	Name string
}

// Launcher defines the interface each launcher must satisfy.
type SourceChange struct {
	// IsAdd is true when the source is being added, and false when it is being removed.
	IsAdd bool

	// Source is the source being added or removed.
	Source *LogSource
}

// Subscription is the type that other components should provide in order to
// subscribe to SourceChanges.
type Subscription = subscriptions.Subscription[SourceChange]

// Subscribe creates a new subscription to this component.
func Subscribe() Subscription {
	return subscriptions.NewSubscription[SourceChange]()
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	componentName,
	fx.Provide(newSourceMgr),
)
