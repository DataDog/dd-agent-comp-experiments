// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package flare

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// registration is provided by other components in order to register a callback
// that will create files in a flare.
type registration struct {
	// callback is called to create the file(s) within a temporary directory.
	callback func(flareDir string) error
}

// CallbackRegistration creates a Registration that will call the given
// callback with a directory in which it should create the necessary files.
// The callback may be called concurrently with any other activity.
func CallbackRegistration(callback func(flareDir string) error) Registration {
	return Registration{
		Registration: registration{
			callback: callback,
		},
	}
}

// FileRegistration creates a Registration that will generate a single file of
// the given name, with the content returned by `callback`.  The callback may be called
// concurrently with any other activity.
func FileRegistration(filename string, callback func() (string, error)) Registration {
	reg := registration{
		callback: func(flareDir string) error {
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
	return Registration{Registration: reg}
}
