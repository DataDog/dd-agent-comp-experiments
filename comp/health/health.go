// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/flare"
	"github.com/djmitche/dd-agent-comp-experiments/comp/ipcapi"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
)

type health struct {
	// Mutex covers all fields, including all componentHealth values
	sync.Mutex

	// started is true once the component has started
	started bool

	// disabled indicates that the component should do nothing.
	disabled bool

	// components maps component package path to that component's current health status
	components map[string]ComponentHealth

	// actorRegistrations stores the set of actorRegistration instances that must
	// be started when the component starts.
	actorRegistrations []*ActorRegistration

	// log supports logging about changes in health status
	log log.Component

	// ipcapi is used in GetHealthRemote
	ipcapi ipcapi.Component
}

type dependencies struct {
	fx.In
	Lc     fx.Lifecycle
	Params ModuleParams `optional:"true"`
	Log    log.Component
	Flare  flare.Component
	IpcAPI ipcapi.Component
}

func newHealth(deps dependencies) Component {
	h := &health{
		disabled:   deps.Params.Disabled,
		components: make(map[string]ComponentHealth),
		log:        deps.Log,
		ipcapi:     deps.IpcAPI,
	}
	deps.Lc.Append(fx.Hook{OnStart: h.start})
	deps.IpcAPI.Register("/agent/health", h.ipcHandler)
	deps.Flare.RegisterFile("health.json", h.flareFile)
	return h
}

// RegisterSimple implements Component#RegisterSimple.
func (h *health) RegisterSimple(component string) *SimpleRegistration {
	h.Lock()
	defer h.Unlock()

	if h.started {
		panic("Health component has already started")
	}

	if !h.disabled {
		if _, exists := h.components[component]; exists {
			panic(fmt.Sprintf("Component %s is already registered with the health component", component))
		}
		h.components[component] = ComponentHealth{Healthy: true}
	}
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

	if !h.disabled {
		if _, exists := h.components[component]; exists {
			panic(fmt.Sprintf("Component %s is already registered with the health component", component))
		}
		h.components[component] = ComponentHealth{Healthy: true}
	}
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

// GetHealthRemote implements Component#GetHealthRemote.
func (h *health) GetHealthRemote() (map[string]ComponentHealth, error) {
	var content map[string]ComponentHealth
	err := h.ipcapi.GetJSON("/agent/health", &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// start starts all actor registrations.
func (h *health) start(ctx context.Context) error {
	h.Lock()
	defer h.Unlock()

	h.started = true
	if !h.disabled {
		for _, reg := range h.actorRegistrations {
			reg.start()
		}
		h.actorRegistrations = nil
	}
	return nil
}

// ipcHandler serves the /agent/health endpoint
func (h *health) ipcHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header()["Content-Type"] = []string{"application/json; charset=UTF-8"}
	json.NewEncoder(w).Encode(h.GetHealth())
}

// flareFile creates the health.json file for Agent flares.
func (h *health) flareFile() (string, error) {
	var bldr strings.Builder
	json.NewEncoder(&bldr).Encode(h.GetHealth())
	return bldr.String(), nil
}

// setHealth sets the health of a specific component.  It is called from the
// XxxRegistration types.
func (h *health) setHealth(component string, healthy bool, message string) {
	h.Lock()
	defer h.Unlock()

	if ch, found := h.components[component]; found && !h.disabled {
		// XXX: we will probably want to do more than just log
		if healthy && !ch.Healthy {
			h.log.Debug(fmt.Sprintf("Component %s is now healthy", component))
		}
		if !healthy && ch.Healthy {
			h.log.Debug(fmt.Sprintf("Component %s is now unhealthy: %s", component, message))
		}
		h.components[component] = ComponentHealth{
			Healthy: healthy,
			Message: message,
		}
	}
}

type ComponentHealth struct {
	Healthy bool
	Message string
}
