// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

import (
	"testing"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestSimple(t *testing.T) {
	var h Component
	app := fxtest.New(t,
		Module,
		log.Module,
		fx.Populate(&h),
	)
	reg := h.RegisterSimple("comp/thing")
	defer app.RequireStart().RequireStop()

	require.Equal(t, ComponentHealth{healthy: true}, h.GetHealth()["comp/thing"])
	reg.SetUnhealthy("uhoh")
	require.Equal(t, ComponentHealth{healthy: false, message: "uhoh"}, h.GetHealth()["comp/thing"])
	reg.SetHealthy()
	require.Equal(t, ComponentHealth{healthy: true}, h.GetHealth()["comp/thing"])
}

func TestActor(t *testing.T) {
	var h Component
	app := fxtest.New(t,
		Module,
		log.Module,
		fx.Populate(&h),
	)
	reg := h.RegisterActor("comp/thing", time.Millisecond)
	defer app.RequireStart().RequireStop()

	for i := 0; i < 3; i++ {
		// signal health
		<-reg.Chan()
		require.Equal(t, ComponentHealth{healthy: true}, h.GetHealth()["comp/thing"])
	}

	// fail to check the messages..
	time.Sleep(10 * time.Millisecond)
	require.Equal(t,
		ComponentHealth{healthy: false, message: "health check timed out"},
		h.GetHealth()["comp/thing"])

	// have to read at least two messages in time to be considered healthy again..
	for i := 0; i < 2; i++ {
		<-reg.Chan()
	}
	require.Equal(t, ComponentHealth{healthy: true}, h.GetHealth()["comp/thing"])

	// stop monitoring
	reg.Stop()

	// fail to check the messages..
	time.Sleep(5 * time.Millisecond)
	require.Equal(t, ComponentHealth{healthy: true}, h.GetHealth()["comp/thing"])
}
