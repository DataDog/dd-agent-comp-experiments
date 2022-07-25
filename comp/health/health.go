// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
)

type health struct {
	// Mutex covers all fields, including all componentHealth values
	sync.Mutex

	// started is true once the component has started
	started bool

	// components maps component package path to that component's current health status
	components map[string]ComponentHealth

	// actorRegistrations stores the set of actorRegistration instances that must
	// be started when the component starts.
	actorRegistrations []*ActorRegistration

	// log supports logging about changes in health status
	log log.Component
}

func newHealth(lc fx.Lifecycle, log log.Component) Component {
	h := &health{
		components: make(map[string]ComponentHealth),
		log:        log,
	}
	lc.Append(fx.Hook{OnStart: h.start})
	return h
}

// RegisterSimple implements Component#RegisterSimple.
func (h *health) RegisterSimple(component string) *SimpleRegistration {
	h.Lock()
	defer h.Unlock()

	if h.started {
		panic("Health component has already started")
	}

	if _, exists := h.components[component]; exists {
		panic(fmt.Sprintf("Component %s is already registered with the health component", component))
	}
	h.components[component] = ComponentHealth{healthy: true}
	return &SimpleRegistration{
		health:    h,
		component: component,
	}
}

// RegisterActor implements Component#RegisterActor.
func (h *health) RegisterActor(component string, healthDuration time.Duration) *ActorRegistration {
	h.Lock()
	defer h.Unlock()

	if h.started {
		panic("Health component has already started")
	}

	if _, exists := h.components[component]; exists {
		panic(fmt.Sprintf("Component %s is already registered with the health component", component))
	}
	h.components[component] = ComponentHealth{healthy: true}
	reg := &ActorRegistration{
		SimpleRegistration: SimpleRegistration{
			health:    h,
			component: component,
		},
		duration:   healthDuration,
		healthChan: make(chan struct{}, 1), // capacity=1 to allow one tick to elapse before failing
		stopped:    make(chan struct{}),
	}
	h.actorRegistrations = append(h.actorRegistrations, reg)
	return reg
}

// GetHealth implements Component#GetHealth.
func (h *health) GetHealth() map[string]ComponentHealth {
	h.Lock()
	defer h.Unlock()

	rv := map[string]ComponentHealth{}
	for k, v := range h.components {
		rv[k] = v
	}
	return rv
}

// start starts all actor registrations.
func (h *health) start(ctx context.Context) error {
	h.Lock()
	defer h.Unlock()

	h.started = true
	for _, reg := range h.actorRegistrations {
		reg.start()
	}
	h.actorRegistrations = nil
	return nil
}

// setHealth sets the health of a specific component.  It is called from the
// XxxRegistration types.
func (h *health) setHealth(component string, healthy bool, message string) {
	h.Lock()
	defer h.Unlock()

	if ch, found := h.components[component]; found {
		// XXX: we will probably want to do more than just log
		if healthy && !ch.healthy {
			h.log.Debug(fmt.Sprintf("Component %s is now healthy", component))
		}
		if !healthy && ch.healthy {
			h.log.Debug(fmt.Sprintf("Component %s is now unhealthy: %s", component, message))
		}
		h.components[component] = ComponentHealth{
			healthy: healthy,
			message: message,
		}
	}
}

type ComponentHealth struct {
	healthy bool
	message string
}
