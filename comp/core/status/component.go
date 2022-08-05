// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package status implements the functionality behind `agent status`.
//
// The data included in the status output is provided by other components, by providing a
// *status.Registration instance in value-group "health".  Nil *status.Registrations will
// be ignored, assuming they are for disabled components.
//
// All of the component's methods can be called concurrently.
package status

import (
	"go.uber.org/fx"
)

// team: agent-shared-components

const componentName = "comp/core/status"

// Component is the component type.
type Component interface {
	// GetStatus gets the agent status.  If the section parameter is not empty, then
	// only that section's status is returned.  This is a newline-terminated string.
	GetStatus(section string) string

	// GetStatus gets the same value as GetStatus, but using the IPC API.
	GetStatusRemote(section string) (string, error)
}

// Registration is provided by other components to register themselves to
// provide status.
type Registration struct {
	fx.Out

	Registration registration `group:"true"`
}

// NewRegistration creates a new Registration.
//
// The section name allows users to select a single section for output (`agent
// status <section-name>`). When all sections are included, they are ordered by
// `order`.  The `cb` returns the text of the section, including the header. If
// an error occurs in `cb`, it should include the error message in its output.
func NewRegistration(section string, order int, cb func() string) Registration {
	return Registration{
		Registration: registration{section, order, cb},
	}
}

// Module defines the fx options for this component.
var Module = fx.Module(
	componentName,
	fx.Provide(newStatus),
)
