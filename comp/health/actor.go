// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

import "time"

// ActorRegistration is the result of registering with RegisterActor.
type ActorRegistration struct {
	SimpleRegistration

	// duration is the maximum time allowed for the component to respond.
	duration time.Duration

	// healthChan is a channel from which the actor must read within the
	// configured duration, or be considered unhealthy.
	healthChan chan struct{}

	// stopped is closed when the monitored component has stopped and thus
	// should no longer be monitored.
	stopped chan struct{}
}

// start begins monitoring of this component
func (reg *ActorRegistration) start() {
	go func() {
		tkr := time.NewTicker(reg.duration)
		for {
			select {
			case <-reg.stopped:
				return
			case <-tkr.C:
				select {
				case <-reg.stopped:
					return
				case reg.healthChan <- struct{}{}:
					reg.SetHealthy()
				default:
					reg.SetUnhealthy("health check timed out")
				}
			}
		}
	}()
}

// Chan returns a channel from which the actor must read within the configured
// duration, or be considered unhealthy.  This is the same channel on every
// call.
func (reg *ActorRegistration) Chan() <-chan struct{} {
	return reg.healthChan
}

// Stop permanently stops monitoring of this component.  Call this method when the
// monitored component stops, to avoid spurious check failures.  This function must
// only be called once.
func (reg *ActorRegistration) Stop() {
	close(reg.stopped)
}
