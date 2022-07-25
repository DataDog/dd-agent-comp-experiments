// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package manager

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/fx"
)

type manager struct {
	// Mutex protects all content in this struct
	sync.Mutex

	// started is true if the manager has been started
	started bool

	// launchers contains the set of registered launchers
	launchers map[string]Launcher
}

func newManager(lc fx.Lifecycle) Component {
	m := &manager{
		started:   false,
		launchers: make(map[string]Launcher),
	}
	lc.Append(fx.Hook{OnStart: m.start})
	return m
}

// start updates the manager's 'started' flag, affecting which method calls are
// valid.
//
// Note that, because the individual launchers depend on this component, they may
// not have been started yet!
func (m *manager) start(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	m.started = true
	return nil
}

// RegisterLauncher implements Component#RegisterLauncher.
func (m *manager) RegisterLauncher(name string, launcher Launcher) error {
	m.Lock()
	defer m.Unlock()

	if m.started {
		panic(fmt.Sprintf("Launcher %s cannot be registered after component startup", name))
	}

	if _, exists := m.launchers[name]; exists {
		return fmt.Errorf("Launcher %s has already been registered", name)
	}

	m.launchers[name] = launcher
	return nil
}

// GetLaunchers implements Component#GetLaunchers.
func (m *manager) GetLaunchers() map[string]Launcher {
	m.Lock()
	defer m.Unlock()

	if !m.started {
		panic("GetLaunchers cannot be called until the manager has started")
	}

	return m.launchers
}

// GetLauncher implements Component#GetLauncher.
func (m *manager) GetLauncher(name string) Launcher {
	m.Lock()
	defer m.Unlock()

	if !m.started {
		panic("GetLauncher cannot be called until the manager has started")
	}

	return m.launchers[name]
}
