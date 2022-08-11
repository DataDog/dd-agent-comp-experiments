// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package actor

import (
	"context"
	"time"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/health"
)

// MonitorLiveness indicates that the actor should report its "liveness" to the
// given Health handle.
//
// The given period should be comfortably longer than the longest time between
// runs of the component's main loop.
func (a *Actor) MonitorLiveness(handle *health.Handle, period time.Duration) {
	a.healthHandle = handle
	a.livenessPeriod = period
}

// This method must not be called before the monitored component has started.
func (a *Actor) livenessMonitor() (<-chan struct{}, func()) {
	// if MonitorLiveness wasn't called, do nothing
	if a.healthHandle == nil {
		return make(chan struct{}), func() {}
	}

	ch := make(chan struct{}, 1) // capacity=1 to allow one tick to elapse before failure
	stopped := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer close(stopped)
		tkr := time.NewTicker(a.livenessPeriod)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tkr.C:
				select {
				case <-ctx.Done():
					return
				case ch <- struct{}{}:
					// we were able to add an item to the channel, so the
					// component is healthy
					a.healthHandle.SetHealthy()
				default:
					// we are not stopped, and were not able to write an item,
					// so the component is unhealthy
					a.healthHandle.SetUnhealthy("health check timed out")
				}
			}
		}
	}()

	stop := func() {
		cancel()
		<-stopped // wait until goroutine has stopped
	}

	return ch, stop
}
