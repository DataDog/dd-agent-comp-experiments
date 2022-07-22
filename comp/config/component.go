// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package config implements a component to handle agent configuration.  This
// component wraps Viper.
//
// The component loads configuration when it is initialized, so that other
// dependent components can use it during their setup phase.  As such, it needs
// the config file path passed to its Module function.
package config

import "go.uber.org/fx"

// Component is the component type.
type Component interface {
	// GetInt gets an integer-typed config parameter value.
	GetInt(key string) int

	// GetBool gets an integer-typed config parameter value.
	GetBool(key string) bool

	// GetInt gets a string-typed config parameter value.
	GetString(key string) string
}

// Module defines the fx options for this component.
func Module(configFilePath string) fx.Option {
	return fx.Provide(func() (Component, error) {
		return newConfig(configFilePath)
	})
}
