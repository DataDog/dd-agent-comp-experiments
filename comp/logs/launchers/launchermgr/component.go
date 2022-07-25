// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package launchermgr implements a component managing logs-agent launchers.  It collects
// the set of loaded launchers during start-up, and allows enumeration and retrieval
// as necessary.
//
// Launchers register themselves with this component.
//
// All component methods can be called concurrently.
package launchermgr

import "go.uber.org/fx"

// Component is the component type.
type Component interface {
	// RegisterLaucher registers a launcher with the manager.  This must be called
	// before the manager is started.  It is an error to register multiple launchers
	// with the same name.
	RegisterLauncher(name string, launcher Launcher) error

	// GetLaunchers gets a map of launchers by name.  This method must be
	// called after the manager has started, guaranteeing that the map is
	// immutable.  Callers must not modify the map.
	GetLaunchers() map[string]Launcher

	// GetLauncher gets a launcher by name, or nil if no such launcher exists.
	// This method must be called after the manager has started.
	GetLauncher(name string) Launcher
}

// Launcher defines the interface each launcher must satisfy.
type Launcher interface {
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	"comp/logs/launchers/manager",
	fx.Provide(newManager),
)
