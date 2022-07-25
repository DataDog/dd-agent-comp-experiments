// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package config implements a component to handle agent configuration.  This
// component wraps Viper.
//
// The component's Setup method must be called before any other components which
// might call is other methods.  This is typically accomplished by calling Setup
// from an application-level `fx.Invoke`.
package config

import "go.uber.org/fx"

// Component is the component type.
type Component interface {
	// Setup sets up the component.  It must be called before any of the other
	// methods.
	Setup(configFilePath string)

	// GetInt gets an integer-typed config parameter value.
	GetInt(key string) int

	// GetBool gets an integer-typed config parameter value.
	GetBool(key string) bool

	// GetInt gets a string-typed config parameter value.
	GetString(key string) string
}

// Module defines the fx options for this component.
var Module = fx.Module(
	"comp/config",
	fx.Provide(newConfig),
)
