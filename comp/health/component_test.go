// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

import (
	"testing"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestSimple(t *testing.T) {
	type provides struct {
		fx.Out
		Registration *Registration `group:"true"`
	}

	var h Component
	reg := &Registration{component: "comp/thing"}
	app := fxtest.New(t,
		Module,
		config.Module,
		log.Module,
		fx.Provide(func() provides {
			return provides{
				Registration: reg,
			}
		}),
		fx.Populate(&h),
	)
	defer app.RequireStart().RequireStop()

	require.Equal(t, ComponentHealth{Healthy: true}, h.GetHealth()["comp/thing"])
	reg.SetUnhealthy("uhoh")
	require.Equal(t, ComponentHealth{Healthy: false, Message: "uhoh"}, h.GetHealth()["comp/thing"])
	reg.SetHealthy()
	require.Equal(t, ComponentHealth{Healthy: true}, h.GetHealth()["comp/thing"])
}

func TestLiveness(t *testing.T) {
	type provides struct {
		fx.Out
		Registration *Registration `group:"true"`
	}

	var h Component
	reg := &Registration{component: "comp/thing"}
	app := fxtest.New(t,
		Module,
		config.Module,
		log.Module,
		fx.Provide(func() provides {
			return provides{
				Registration: reg,
			}
		}),
		fx.Populate(&h),
	)
	defer app.RequireStart().RequireStop()

	ch, stop := reg.LivenessMonitor(time.Millisecond)

	for i := 0; i < 3; i++ {
		// signal health
		<-ch
		require.Equal(t, ComponentHealth{Healthy: true}, h.GetHealth()["comp/thing"])
	}

	// fail to check the messages..
	time.Sleep(10 * time.Millisecond)
	require.Equal(t,
		ComponentHealth{Healthy: false, Message: "health check timed out"},
		h.GetHealth()["comp/thing"])

	// have to read at least two messages in time to be considered healthy again..
	for i := 0; i < 2; i++ {
		<-ch
	}
	require.Equal(t, ComponentHealth{Healthy: true}, h.GetHealth()["comp/thing"])

	// stop monitoring
	stop()

	// fail to check the messages..
	time.Sleep(5 * time.Millisecond)
	require.Equal(t, ComponentHealth{Healthy: true}, h.GetHealth()["comp/thing"])
}
