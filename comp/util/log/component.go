// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package log implements a component to handle logging internal to the agent.
//
// The log methods can be called at any point in the component's lifecycle, but
// will be buffered until the component starts and only written after that
// time.
package log

import "go.uber.org/fx"

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// Configure defines the settings for the logger.  This can be called
	// before the component starts, such as in an fx.Invoke.  It must only be
	// called once, typically from the App initialization.
	Configure(level string) error

	// Debug logs at the debug level.
	Debug(v ...interface{})

	// Flush flushes the underlying inner log
	Flush()

	// ..more methods, obviously :)
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Provide(newLogger)
