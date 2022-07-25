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
// The health component is ready for registration as soon as it is initialized, and
// registration must be completed before the health component starts.
//
// All of the component's methods can be called concurrently.
package health

import (
	"time"

	"go.uber.org/fx"
)

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// RegisterSimple registers a component for "simple" monitoring.  It is assumed
	// to be healthy initially, and that status can be updated with methods on the
	// returned value.
	//
	// Component is the component's package path (e.g., `comp/health`).
	RegisterSimple(component string) *SimpleRegistration

	// RegisterActor register a component for "actor" monitoring.  Once the app
	// starts, the actor must read from the channel in the returned value
	// within the given duration, or it will be considered unhealthy.
	//
	// To use: call `RegisterActor` in the component constructor, and store the
	// resulting registration.  At the beginning of the actor's run function,
	// defer a call to reg.Stop().  Within the actor's run loop, read from
	// reg.Chan().
	//
	// Component is the component's package path (e.g., `comp/health`).
	RegisterActor(component string, healthDuration time.Duration) *ActorRegistration

	// GetHealth gets a map containing the health of all components.  This map is a copy
	// and will not be altered after return.
	GetHealth() map[string]ComponentHealth
}

// Module defines the fx options for this component.
var Module = fx.Module(
	"comp/health",
	fx.Provide(newHealth),
)
