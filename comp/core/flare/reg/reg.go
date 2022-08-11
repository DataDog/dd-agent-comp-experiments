// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package reg contains the Registration type and helper functions for
// comp/core/flare, isolated here to prevent Go package cycles.
//
// In most cases, these items can be used from the comp/core/flare package directly.
package reg

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Registration is provided by other components in order to register a callback
// that will create files in a flare.
type Registration struct {
	// Callback is called to create the file(s) within a temporary directory.
	Callback func(flareDir string) error
}

// FileRegistration creates a Registration that will generate a single file of
// the given name, with the content returned by `callback`.  The callback may be called
// concurrently with any other activity.
func FileRegistration(filename string, callback func() (string, error)) *Registration {
	return &Registration{
		Callback: func(flareDir string) error {
			content, err := callback()
			if err != nil {
				return err
			}

			fullpath := filepath.Join(flareDir, filename)
			parentDir := filepath.Dir(fullpath)
			err = os.MkdirAll(parentDir, 0o700)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(fullpath, []byte(content), 0o600)
		},
	}
}
