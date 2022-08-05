// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package health implements a component that monitors the health of other
// components.
//
// The data from this component is provided by other components, by providing a
// health.Registration instance in value-group "health".
//
// The health component supports monitoring several kinds of components.  All begin
// by providing a Registration in their constructor.  The constructor should capture
// and store the Handle within this Registration for later use.
//
// "Simple" components record their health status with this component with
// simple function calls on the Handle.  If such a component deadlocks, the
// health component will be unaware.
//
// Actor-based components can register and receive a channel from which they must
// read within a configured amount of time.  This approach makes sense for components
// using the [actor model](https://en.wikipedia.org/wiki/Actor_model), where the
// component is considered unhealthy if it is not polling for events frequently.
//
// All of the component's methods can be called concurrently.
package health

import (
	"go.uber.org/fx"
)

// team: agent-shared-components

const componentName = "comp/core/health"

// Component is the component type.
type Component interface {
	// GetHealth gets a map containing the health of all components.  This map is a copy
	// and will not be altered after return.
	GetHealth() map[string]ComponentHealth

	// GetHealthRemote gets the same value as GetHealth, but using the IPC API.
	GetHealthRemote() (map[string]ComponentHealth, error)
}

// Registration is provided by other components to register themselves to
// provide flare data.
type Registration struct {
	fx.Out

	Handle *Handle `group:"true"`
}

// NewRegistration creates a new Registration instance for the named component.
func NewRegistration(component string) Registration {
	return Registration{
		Handle: &Handle{component: component},
	}
}

// Module defines the fx options for this component.
var Module = fx.Module(
	componentName,
	fx.Provide(newHealth),
)
