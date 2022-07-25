// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package file implements a manager.Launcher that responds to file sources by starting
// tailers for the indicated files.
package file

import (
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/manager"
	"go.uber.org/fx"
)

// Component is the component type.
type Component interface {
	// Launcher includes the manager's launcher methods here
	manager.Launcher
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	"comp/logs/launchers/file",
	fx.Provide(newLauncher),
)
