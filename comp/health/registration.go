// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

import (
	"context"
	"time"
)

// Registration represents a registration with this component.
type Registration struct {
	health    *health
	component string
}

// SetUnhealthy records this component as being unhealthy, with the included message
// summarizing the problem.  This can be called at any time.
func (reg *Registration) SetUnhealthy(message string) {
	reg.health.setHealth(reg.component, false, message)
}

// SetHealthy records this component as being healthy.  This can be called at any time.
func (reg *Registration) SetHealthy() {
	reg.health.setHealth(reg.component, true, "")
}

// LivenessMonitor starts a goroutine that periodically writes items to the
// returned channel.  The expectation is that the monitored component will read
// from this channel in its main loop.  If the channel fills up, the component
// is presumed to be broken or stuck, and marked unhealthy.
//
// The given period should be comfortably longer than the longest time between
// runs of the component's main loop.
//
// The returned function cancels the liveness monitor, leaving the component in its
// existing state.
func (reg *Registration) LivenessMonitor(period time.Duration) (<-chan struct{}, func()) {
	ch := make(chan struct{}, 1) // capacity=1 to allow one tick to elapse before failure
	stopped := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer close(stopped)
		tkr := time.NewTicker(period)
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
					reg.SetHealthy()
				default:
					// we are not stopped, and were not able to write an item,
					// so the component is unhealthy
					reg.SetUnhealthy("health check timed out")
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
