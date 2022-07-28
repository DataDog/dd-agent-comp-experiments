// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package status implements the functionality behind `agent status`.
//
// It allows other components to register "sections" of the status output, and makes
// these available via the ipcapi.  The `agent status` command can then choose which
// sections to display.
//
// All of the component's methods can be called concurrently.
package status

import (
	"go.uber.org/fx"
)

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// RegisterSection registers a callback that will get the text of a
	// section.  The order parameter determines the order in which sections are
	// displayed when displaying all sections.  The callback must produce some
	// output; if an error occurs, it should include that error in the output.
	RegisterSection(section string, order int, cb func() string)

	// GetStatus gets the agent status.  If the section parameter is not empty, then
	// only that section's status is returned.  This is a newline-terminated string.
	GetStatus(section string) string

	// GetStatus gets the same value as GetStatus, but using the IPC API.
	GetStatusRemote(section string) (string, error)
}

// Module defines the fx options for this component.
var Module = fx.Module(
	"comp/status",
	fx.Provide(newStatus),
)
