// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package flare implements a component creates flares for submission to support.
//
// Other components depend on this one, registering callbacks that can be used
// to create content in the flare.
//
// This component registers itself with the ipcapi component, and supports either
// generating a flare locally (CreateFlare) or calling the API to direct the running
// Agent to create a flare (CreateFlareRemote).  Creating a flare locally in a
// process that is not running a full Agent will still generate a flare, but that
// flare will lack information from components that are not running.
//
// All flare methods can be called at any time, but the expectation is that Register*
// methods would be called during the setup phase.
package flare

// NOTE: it might be nice to users to generate a "README.md" describing each file in
// the flare, based on docs passed to the flare.Register* methods.

import (
	"testing"

	"go.uber.org/fx"
)

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// RegisterFile registers a callback that will create content to place into
	// a file in the flare.  The returned content will be scrubbed of secrets.
	// If the callback fails, the error will be recorded in FLARE-ERRORS.txt,
	// but flare generation will continue.  The given filename must be
	// relative.
	RegisterFile(filename string, cb func() (string, error))

	// Register registers a callback that can write whatever data it wants
	// to the flare, by creating files and directories under flareDir.
	//
	// The callback must scrub any content it writes to the flareDir.
	Register(cb func(flareDir string) error)

	// CreateFlare creates a new flare locally and returns the path to the
	// flare file.
	CreateFlare() (string, error)

	// CreateFlareRemote calls the running Agent's IPC API to instruct it to
	// generate a flare remotely.
	CreateFlareRemote() (string, error)
}

// Mock implements mock-specific methods.
type Mock interface {
	Component

	// GetFlareFile generates a temporary flare, then returns the content of the
	// named file.
	GetFlareFile(t *testing.T, filename string) (string, error)
}

// Module defines the fx options for this component.
var Module = fx.Module(
	"comp/flare",
	fx.Provide(newFlare),
)

// MockModule defines the fx options for the mock component.
var MockModule = fx.Module(
	"comp/flare",
	fx.Provide(newMock),
)
