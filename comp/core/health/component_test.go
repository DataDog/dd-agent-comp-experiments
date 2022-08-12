// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

import (
	"testing"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/startup"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestSimple(t *testing.T) {
	var h Component
	reg := NewRegistration("comp/thing")
	app := fxtest.New(t,
		Module,
		log.Module,
		config.MockModule,
		fx.Supply(internal.BundleParams{AutoStart: startup.Always}),
		fx.Supply(reg),
		fx.Populate(&h),
	)
	defer app.RequireStart().RequireStop()

	require.Equal(t, ComponentHealth{Healthy: true}, h.GetHealth()["comp/thing"])
	reg.Handle.SetUnhealthy("uhoh")
	require.Equal(t, ComponentHealth{Healthy: false, Message: "uhoh"}, h.GetHealth()["comp/thing"])
	reg.Handle.SetHealthy()
	require.Equal(t, ComponentHealth{Healthy: true}, h.GetHealth()["comp/thing"])
}
