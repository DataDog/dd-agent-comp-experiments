// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package file

import (
	"testing"

	"github.com/djmitche/dd-agent-comp-experiments/comp/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal/sourcemgr"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/launchermgr"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestMyComponent(t *testing.T) {
	var comp Component
	var smgr sourcemgr.Component
	app := fxtest.New(t,
		Module,
		health.Module,
		log.Module,
		sourcemgr.Module,
		launchermgr.Module,
		fx.Populate(&comp),
		fx.Populate(&smgr),
	)

	defer app.RequireStart().RequireStop()

	smgr.AddSource(&sourcemgr.LogSource{Name: "testy"})
	// if this launcher actually launched anything, we'd assert on that here
	// with require.Eventually(t, ...)
}
