// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package status

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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
	sections []sectionInfo

	// ipcapi is used to serve the status to remote instances.
	ipcapi ipcapi.Component
}

type dependencies struct {
	fx.In
	Lc     fx.Lifecycle
	Flare  flare.Component
	IpcAPI ipcapi.Component
}

func newStatus(deps dependencies) Component {
	s := &status{
		ipcapi: deps.IpcAPI,
	}
	deps.IpcAPI.Register("/agent/status", s.ipcHandler)
	deps.Flare.RegisterFile("agent-status.txt", s.flareFile)
	return s
}

// RegisterSection implements Component#RegisterSection.
func (s *status) RegisterSection(section string, order int, cb func() string) {
	s.Lock()
	defer s.Unlock()

	for _, s := range s.sections {
		if s.name == section {
			panic(fmt.Sprintf("Section %s is already registered with the status component", section))
		}
	}

	s.sections = append(s.sections, sectionInfo{name: section, order: order, cb: cb})
	sort.Sort(byOrder(s.sections))
}

// GetStatus implements Component#GetStatus.
func (s *status) GetStatus(section string) string {
	s.Lock()
	defer s.Unlock()

	var bldr strings.Builder
	for _, s := range s.sections {
		if section != "" && s.name != section {
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
	var content map[string]string
	// TODO: support getting just one section

	err := s.ipcapi.GetJSON("/agent/status", &content)
	if err != nil {
		return "", err
	}

	return content["status"], nil
}

// ipcHandler serves the /agent/status endpoint
func (s *status) ipcHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header()["Content-Type"] = []string{"application/json; charset=UTF-8"}
	json.NewEncoder(w).Encode(map[string]string{"status": s.GetStatus("")})
}

// flareFile creates the agent-status.txt file for flares.
func (s *status) flareFile() (string, error) {
	return s.GetStatus(""), nil
}

type sectionInfo struct {
	name  string
	order int
	cb    func() string
}

// byOrder supports sorting sections by order.
type byOrder []sectionInfo

func (a byOrder) Len() int           { return len(a) }
func (a byOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byOrder) Less(i, j int) bool { return a[i].order < a[j].order }
