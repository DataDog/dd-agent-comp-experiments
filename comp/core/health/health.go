// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	flare "github.com/DataDog/dd-agent-comp-experiments/comp/core/flare"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcserver"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"go.uber.org/fx"
)

type health struct {
	// Mutex covers all fields, including all componentHealth values
	sync.Mutex

	// components maps component package path to that component's current health status
	components map[string]ComponentHealth

	// log supports logging about changes in health status
	log log.Component
}

type dependencies struct {
	fx.In

	Lc     fx.Lifecycle
	Params internal.BundleParams
	Log    log.Component

	Handles []*Handle `group:"health"`
}

type provides struct {
	fx.Out

	Component
	FlareReg flare.Registration
	IPCRoute ipcserver.Route
}

func newHealth(deps dependencies) provides {
	h := &health{
		components: make(map[string]ComponentHealth),
		log:        deps.Log,
	}

	// provide each registration with a pointer to the new component, and
	// default to a healthy status. The Handles will update the component
	// as health status changes.
	for _, handle := range deps.Handles {
		handle.health = h
		h.components[handle.component] = ComponentHealth{Healthy: true}
	}

	return provides{
		Component: h,
		FlareReg:  flare.FileRegistration("health.json", h.flareFile),
		IPCRoute:  ipcserver.NewRoute("/agent/health", h.ipcHandler),
	}
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

	if ch, found := h.components[component]; found {
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
