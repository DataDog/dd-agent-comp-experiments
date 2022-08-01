// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package log implements a component to handle logging internal to the agent.
//
// The log methods can be called at any point in the component's lifecycle, but
// will be buffered until the component starts and only written after that
// time.
//
// This component has special support for acting as an Fx logger, including
// capturing log messages emitted before the component itself is initialized.
package log

import (
	"testing"

	"go.uber.org/fx"
)

// team: agent-shared-components

const componentName = "comp/util/log"

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

// Mock is the mocked component type.
type Mock interface {
	Component

	// SetT sets the testing instance to which log messages should be copied.
	// If this is not called and ModuleParams#Console is false, log messages
	// are not written anywhere.
	SetT(*testing.T)

	// StartCapture begins capturing log messages.  All log messages are
	// captured, regardless of level.
	StartCapture()

	// Captured returns the log messages captured so far.  The returned slice
	// is a copy and will not be modified after return
	Captured() []string

	// EndCapture ends capturing log messages and discards buffered log
	// messages.  It's not required to call this.
	EndCapture()
}

// ModuleParams are the parameters to Module.
type ModuleParams struct {
	// Console determines whether log messages should be output to the console.
	Console bool
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	componentName,
	fx.Provide(newLogger),
)

// MockModule defines the fx options for the mock component.
var MockModule fx.Option = fx.Module(
	componentName,
	fx.Provide(newMockLogger),
)
