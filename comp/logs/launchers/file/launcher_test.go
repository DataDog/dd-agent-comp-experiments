// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package file

import (
	"strings"
	"testing"
	"time"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/internal/sourcemgr"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/launchers/launchermgr"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/comptest"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/startup"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestLauncher(t *testing.T) {
	var comp Component
	var smgr sourcemgr.Component
	var l log.Component
	comptest.FxTest(t,
		Module,
		core.MockBundle,
		fx.Supply(internal.BundleParams{AutoStart: startup.Always}),
		sourcemgr.Module,
		launchermgr.Module,
		fx.Populate(&comp),
		fx.Populate(&smgr),
		fx.Populate(&l),
	).WithRunningApp(func() {
		// Arrange
		l.(log.Mock).StartCapture()

		// Act
		smgr.AddSource(&sourcemgr.LogSource{Name: "testy"})

		// Assert
		require.Eventually(t, func() bool {
			for _, m := range l.(log.Mock).Captured() {
				if strings.Contains(m, "got LogSource change") {
					return true
				}
			}
			return false
		}, time.Second, time.Millisecond)
	})
}
