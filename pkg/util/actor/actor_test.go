// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package actor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/health"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/startup"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestWithoutComponents(t *testing.T) {
	actor := Actor{}
	ch := make(chan int)

	run := func(ctx context.Context, alive <-chan struct{}) {
		for {
			select {
			case <-alive:
			case v := <-ch:
				fmt.Printf("GOT: %d\n", v)
			case <-ctx.Done():
				fmt.Println("Stopping")
				return
			}
		}
	}

	actor.Start(run)
	ch <- 1
	ch <- 2
	actor.Stop(context.Background())

	// Output:
	// GOT: 1
	// GOT: 2
	// Stopping
}

type testComp struct {
	actor Actor
}

func newTestComp(lc fx.Lifecycle) (*testComp, health.Registration) {
	reg := health.NewRegistration("test-comp")
	c := &testComp{}
	c.actor.MonitorLiveness(reg.Handle, time.Millisecond)
	c.actor.HookLifecycle(lc, c.run)
	return c, reg
}

func (c *testComp) run(ctx context.Context, alive <-chan struct{}) {
	// this is healthy for about 5ms, then unhealthy for about 5ms, and repeats
	// that pattern.
	tkr := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-alive:
		case <-tkr.C:
			time.Sleep(5 * time.Millisecond) // unhealthy for a few ms
		case <-ctx.Done():
			return
		}
	}
}

func TestWithHealth(t *testing.T) {
	var comp *testComp
	var health health.Component
	app := fxtest.New(t,
		fx.Supply(core.BundleParams{AutoStart: startup.Never}),
		core.Bundle,
		fx.Provide(newTestComp),
		fx.Populate(&comp),
		fx.Populate(&health),
	)

	defer app.RequireStart().RequireStop()

	// see it go unhealthy..
	require.Eventually(t, func() bool {
		return !health.GetHealth()["test-comp"].Healthy
	}, time.Second, time.Millisecond)

	// see it return to healthy..
	require.Eventually(t, func() bool {
		return !health.GetHealth()["test-comp"].Healthy
	}, time.Second, time.Millisecond)
}
