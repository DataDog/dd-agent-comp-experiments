// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package flare

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcserver"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/mholt/archiver"
	"go.uber.org/fx"
)

type flare struct {
	// Mutex covers all fields
	sync.Mutex

	// registrations contains all registrations by other components
	registrations []registration

	// log is the log component
	log log.Component
}

type dependencies struct {
	fx.In

	Config        config.Component
	Log           log.Component
	Registrations []registration `group:"flare"`
}

func newFlare(deps dependencies) (Component, ipcserver.Route) {
	f := &flare{
		registrations: providedRegistrations(deps.Registrations),
		log:           deps.Log,
	}

	return f, ipcserver.NewRoute("/agent/flare", f.ipcHandler)
}

type mockDependencies struct {
	fx.In

	Registrations []registration `group:"flare"`
}

func newMock(deps mockDependencies) Component {
	// mock is just like the real thing, but doesn't use ipcserver or config.
	return &flare{
		registrations: providedRegistrations(deps.Registrations),
	}
}

// providedRegistrations skips regsitrations with a nil callback.
func providedRegistrations(registrations []registration) []registration {
	provided := make([]registration, 0, len(registrations))
	for _, r := range registrations {
		if r.callback != nil {
			provided = append(provided, r)
		}
	}
	return provided
}

// CreateFlare implements Component#CreateFlare.
func (f *flare) CreateFlare() (string, error) {
	f.Lock()
	defer f.Unlock()

	// make a root temporary directory
	rootDir, err := ioutil.TempDir("", "experiment-flare-*")
	if err != nil {
		return "", err
	}

	// TODO: use something like comp/hostname for this.
	flareDir := filepath.Join(rootDir, "hostname")

	err = os.MkdirAll(flareDir, 0o700)
	if err != nil {
		return "", err
	}

	// on completion, remove the flareDir, but leave the archiveFile.
	defer os.RemoveAll(flareDir)

	err = f.writeFlareFiles(flareDir, false)
	if err != nil {
		return "", err
	}

	archiveFile := filepath.Join(rootDir, "hostname-timestamp.zip")
	err = f.createArchive(flareDir, archiveFile)
	if err != nil {
		return "", err
	}

	return archiveFile, nil
}

// GetFlareFile implements Mock#GetFlareFile.
func (f *flare) GetFlareFile(t *testing.T, filename string) (string, error) {
	f.Lock()
	defer f.Unlock()

	flareDir := t.TempDir()
	err := f.writeFlareFiles(flareDir, true)
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadFile(filepath.Join(flareDir, filename))
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ipcHandler serves the /agent/flare endpoint.  On success, this returns a 200
// with {"filename": <filename>} giving the local filename of the flare file.
func (f *flare) ipcHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header()["Content-Type"] = []string{"application/json; charset=UTF-8"}

	f.log.Debug("Creating flare for remote request")

	archiveFile, err := f.CreateFlare()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"filename": archiveFile,
	})
}

// writeFlareFiles calls all of the callbacks to write all flare files to disk.
// If returnErrors is true then errors from callbacks are returned immediately
// (for testing).
//
// XXX note that this does no scrubbing
//
// It assumes f is locked.
func (f *flare) writeFlareFiles(flareDir string, returnErrors bool) error {
	errors := []string{}
	for _, reg := range f.registrations {
		err := reg.callback(flareDir)
		if err != nil {
			if returnErrors {
				return err
			}
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		// attempt to write FLARE-ERRORS.txt; an error here is actually fatal
		err := ioutil.WriteFile(
			filepath.Join(flareDir, "FLARE-ERRORS.txt"),
			[]byte(strings.Join(errors, "\n")),
			0o600)
		if err != nil {
			return err
		}
	}

	return nil
}

// createArchive creates an archive of the flareDir.
//
// It assumes f is locked.
func (f *flare) createArchive(flareDir, archiveFile string) error {
	return archiver.Archive([]string{flareDir}, archiveFile)
}
