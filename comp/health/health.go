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

	flare "github.com/djmitche/dd-agent-comp-experiments/comp/flare"
	"github.com/djmitche/dd-agent-comp-experiments/comp/ipcapi"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
)

type health struct {
	// Mutex covers all fields, including all componentHealth values
	sync.Mutex

	// disabled indicates that the component should do nothing.
	disabled bool

	// components maps component package path to that component's current health status
	components map[string]ComponentHealth

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
	IpcAPI ipcapi.Component
}

type out struct {
	fx.Out

	Component
	FlareReg flare.Registration `group:"flare"`
}

func newHealth(deps dependencies) out {
	h := &health{
		disabled:   deps.Params.Disabled,
		components: make(map[string]ComponentHealth),
		log:        deps.Log,
		ipcapi:     deps.IpcAPI,
	}
	deps.IpcAPI.Register("/agent/health", h.ipcHandler)
	return out{
		Component: h,
		FlareReg:  flare.FileRegistration("health.json", h.flareFile),
	}
}

// RegisterSimple implements Component#RegisterSimple.
func (h *health) Register(component string) *Registration {
	h.Lock()
	defer h.Unlock()

	if !h.disabled {
		if _, exists := h.components[component]; exists {
			panic(fmt.Sprintf("Component %s is already registered with the health component", component))
		}
		h.components[component] = ComponentHealth{Healthy: true}
	}
	return &Registration{
		health:    h,
		component: component,
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

// GetHealthRemote implements Component#GetHealthRemote.
func (h *health) GetHealthRemote() (map[string]ComponentHealth, error) {
	var content map[string]ComponentHealth
	err := h.ipcapi.GetJSON("/agent/health", &content)
	if err != nil {
		return nil, err
	}

	return content, nil
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
