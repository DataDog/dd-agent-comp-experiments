// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package health implements a component that monitors the health of other
// components.
//
// The health component supports monitoring several kinds of components.  "Simple"
// components record their health status with this component with simple function
// calls.  If such a component deadlocks, the health component will be unaware.
//
// Actor-based components can register and receive a channel from which they must
// read within a configured amount of time.  This approach makes sense for components
// using the [actor model](https://en.wikipedia.org/wiki/Actor_model), where the
// component is considered unhealthy if it is not polling for events frequently.
//
// The health component is ready for registration as soon as it is initialized.
// Registration can occur at any time, but typically occurs before components have started.
//
// All of the component's methods can be called concurrently.
package health

import (
	"go.uber.org/fx"
)

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// Register registers a component for monitoring.  It is assumed to be
	// healthy initially, and that status can be updated with methods on the
	// returned value.
	//
	// Component is the component's package path (e.g., `comp/health`).
	Register(component string) *Registration

	// GetHealth gets a map containing the health of all components.  This map is a copy
	// and will not be altered after return.
	GetHealth() map[string]ComponentHealth

	// GetHealthRemote gets the same value as GetHealth, but using the IPC API.
	GetHealthRemote() (map[string]ComponentHealth, error)
}

// ModuleParams are the parameters to Module.
type ModuleParams struct {
	// Disabled indicates that the component should ignore all registration and
	// perform no monitoring.  This is intended for one-shot processes such as
	// `agent status`.
	Disabled bool
}

// Module defines the fx options for this component.
var Module = fx.Module(
	"comp/health",
	fx.Provide(newHealth),
)
