// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package file

import (
	"strings"
	"testing"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal/sourcemgr"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/launchermgr"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestMyComponent(t *testing.T) {
	var comp Component
	var smgr sourcemgr.Component
	var l log.Component
	app := fxtest.New(t,
		Module,
		health.Module,
		config.Module,
		log.MockModule,
		sourcemgr.Module,
		launchermgr.Module,
		fx.Populate(&comp),
		fx.Populate(&smgr),
		fx.Populate(&l),
	)

	l.(log.Mock).SetT(t)
	defer app.RequireStart().RequireStop()

	// Arrange
	l.(log.Mock).StartCapture()

	// Act
	smgr.AddSource(&sourcemgr.LogSource{Name: "testy"})

	// Assert
	require.Eventually(t, func() bool {
		for _, m := range l.(log.Mock).Captured() {
			if strings.Contains(m, "got change") {
				return true
			}
		}
		return false
	}, time.Second, time.Millisecond)
}
