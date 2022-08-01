// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/djmitche/dd-agent-comp-experiments/comp/flare"
	"github.com/djmitche/dd-agent-comp-experiments/comp/ipcapi"
	"go.uber.org/fx"
)

type status struct {
	// Mutex covers all fields, including all componentHealth values
	sync.Mutex

	// components maps component package path to that component's current status
	sections []*Registration
}

type dependencies struct {
	fx.In

	Lc            fx.Lifecycle
	Registrations []*Registration `group:"status"`
}

type provides struct {
	fx.Out

	Component
	FlareReg    *flare.Registration `group:"flare"`
	IpcAPIRoute *ipcapi.Route       `group:"ipcapi"`
}

func newStatus(deps dependencies) provides {
	s := &status{
		sections: deps.Registrations,
	}
	return provides{
		Component:   s,
		FlareReg:    flare.FileRegistration("agent-status.json", s.flareFile),
		IpcAPIRoute: ipcapi.NewRoute("/agent/status", s.ipcHandler),
	}
}

// GetStatus implements Component#GetStatus.
func (s *status) GetStatus(section string) string {
	s.Lock()
	defer s.Unlock()

	var bldr strings.Builder
	for _, s := range s.sections {
		if section != "" && s.section != section {
			continue
		}

		fmt.Fprintf(&bldr, "%s\n", s.cb())
	}

	if bldr.Len() == 0 {
		fmt.Fprintf(&bldr, "Status section %s is not defined", section)
	}

	return bldr.String() + "\n"
}

// GetStatusRemote implements Component#GetStatusRemote.
func (s *status) GetStatusRemote(section string) (string, error) {
	return "", errors.New("TODO")
	/*
		var content map[string]string
		path := "/agent/status"
		if section != "" {
			path = fmt.Sprintf("%s?section=%s", path, section)
		}

		err := s.ipcapi.GetJSON(path, &content)
		if err != nil {
			return "", err
		}

		return content["status"], nil
	*/
}

// ipcHandler serves the /agent/status endpoint
func (s *status) ipcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header()["Content-Type"] = []string{"application/json; charset=UTF-8"}

	var section string
	sections, ok := r.URL.Query()["section"]
	if ok && len(sections) == 1 {
		section = sections[0]
	}

	json.NewEncoder(w).Encode(map[string]string{"status": s.GetStatus(section)})
}

// flareFile creates the agent-status.txt file for flares.
func (s *status) flareFile() (string, error) {
	return s.GetStatus(""), nil
}

// byOrder supports sorting sections by order.
type byOrder []*Registration

func (a byOrder) Len() int           { return len(a) }
func (a byOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byOrder) Less(i, j int) bool { return a[i].order < a[j].order }
