// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package config implements a component to handle agent configuration.  This
// component wraps Viper.
//
// The component loads its configuration immediately on instantiation, so
// configuration is available to all other components, even before they have
// started.  To accomplish this, it requires the config file path in
// its ModuleParams.
//
// The component attempts to load the configuration file at instantiation, failing
// startup if this is not possible.  The mock component does nothing at
// startup, beginning with an empty config.
package config

import "go.uber.org/fx"

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// GetInt gets an integer-typed config parameter value.
	GetInt(key string) int

	// GetBool gets an integer-typed config parameter value.
	GetBool(key string) bool

	// GetInt gets a string-typed config parameter value.
	GetString(key string) string

	// WriteConfig writes the config to the designated file.
	WriteConfig(filename string) error
}

// Mock implements mock-specific methods.
type Mock interface {
	Component

	// TODO: Set..
}

// ModuleParams are the parameters to Module.
type ModuleParams struct {
	// ConfFilePath is the path to the configuration file.
	ConfFilePath string
}

const componentName = "comp/core/config"

// Module defines the fx options for this component.
var Module = fx.Module(
	componentName,
	fx.Provide(newConfig),
)

// MockModule defines the fx options for the mock component.
var MockModule = fx.Module(
	componentName,
	fx.Provide(newMock),
)
