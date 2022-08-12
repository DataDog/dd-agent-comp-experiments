// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package launchermgr implements a component managing logs-agent launchers.  It collects
// the set of loaded launchers during start-up, and allows enumeration and retrieval
// as necessary.
//
// Launchers should provide a launchermgr.Registration instance to register
// themselves with this manager.
//
// All component methods can be called concurrently.
package launchermgr

import "go.uber.org/fx"

// team: agent-metrics-logs

const componentName = "comp/logs/launchers/launchermgr"

// Component is the component type.
type Component interface {
	// GetLaunchers gets a map of launchers by name.  This method must be
	// called after the manager has started, guaranteeing that the map is
	// immutable.  Callers must not modify the map.
	GetLaunchers() map[string]Launcher

	// GetLauncher gets a launcher by name, or nil if no such launcher exists.
	// This method must be called after the manager has started.
	GetLauncher(name string) Launcher
}

// Registration is provided by launchers to register themselves with the manager.
type Registration struct {
	fx.Out

	Registration registration `group:"launchermgr"`
}

// NewRegistration creates a new Registration instance for the named launcher.
func NewRegistration(name string, launcher Launcher) Registration {
	return Registration{
		Registration: registration{name: name, launcher: launcher},
	}
}

// Launcher defines the interface each launcher must satisfy.
type Launcher interface {
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	componentName,
	fx.Provide(newManager),
)
