// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package ipcclient implements a component to access the IPC server remotely.
package ipcclient

import (
	"go.uber.org/fx"
)

// team: agent-shared-components

const componentName = "comp/ipc/ipcclient"

// Component is the component type.
type Component interface {
	// GetJSON gets the body of the server response from the given path, as JSON
	GetJSON(path string, v any) error
}

var Module = fx.Module(
	componentName,
	fx.Provide(newClient),
)
