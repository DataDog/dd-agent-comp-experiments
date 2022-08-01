// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package flare implements a component creates flares for submission to support.
//
// The data from this component is provided by other components, by providing a
// flare.Registration instance in value-group "flare".
//
// This component registers itself with the ipcapi component, and supports either
// generating a flare locally (CreateFlare) or calling the API to direct the running
// Agent to create a flare (CreateFlareRemote).  Creating a flare locally in a
// process that is not running a full Agent will still generate a flare, but that
// flare will lack information from components that are not running.
//
// All flare methods can be called at any time.
package flare

// NOTE: it might be nice to users to generate a "README.md" describing each file in
// the flare, based on the Registrations.

import (
	"testing"

	"go.uber.org/fx"

	"github.com/djmitche/dd-agent-comp-experiments/comp/flare/reg"
)

// team: agent-shared-components

// Component is the component type.
type Component interface {
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

// re-exports (avoiding a Go package dependency loop)

type Registration = reg.Registration

var FileRegistration = reg.FileRegistration

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
