// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package launchermgr

import (
	"go.uber.org/fx"
)

type manager struct {
	// launchers contains the set of registered launchers
	launchers map[string]Launcher
}

type registration struct {
	// name is the name of the launcher
	name string

	// launcher points to the launcher itself
	launcher Launcher
}

type dependencies struct {
	fx.In

	Registrations []registration `group:"launchermgr"`
}

func newManager(deps dependencies) Component {
	m := &manager{
		launchers: make(map[string]Launcher),
	}
	for _, reg := range deps.Registrations {
		m.launchers[reg.name] = reg.launcher
	}
	return m
}

// GetLaunchers implements Component#GetLaunchers.
func (m *manager) GetLaunchers() map[string]Launcher {
	return m.launchers
}

// GetLauncher implements Component#GetLauncher.
func (m *manager) GetLauncher(name string) Launcher {
	return m.launchers[name]
}
