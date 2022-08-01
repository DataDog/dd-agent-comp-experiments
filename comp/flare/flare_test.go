// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package flare

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/mholt/archiver"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFlareMechanics(t *testing.T) {
	flareDir := t.TempDir()

	type provides struct {
		fx.Out

		Registration *Registration `group:"flare"`
	}

	var flare Component
	app := fxtest.New(t,
		Module,
		config.Module,
		fx.Provide(func() provides {
			return provides{
				Registration: FileRegistration("greeting.txt", func() (string, error) {
					return "hello, world", nil
				}),
			}
		}),
		fx.Populate(&flare),
	)

	defer app.RequireStart().RequireStop()

	archiveFile, err := flare.CreateFlare()
	require.NoError(t, err)

	require.NotEqual(t, "", archiveFile)
	err = archiver.Extract(archiveFile, "hostname/greeting.txt", flareDir)
	require.NoError(t, err)

	content, err := ioutil.ReadFile(filepath.Join(flareDir, "hostname", "greeting.txt"))
	require.NoError(t, err)
	require.Equal(t, "hello, world", string(content))

	// this will create a temporary file without t.TempDir, so we must clean it
	// up manually
	os.RemoveAll(archiveFile)
}

func TestMock(t *testing.T) {
	type provides struct {
		fx.Out

		Registration *Registration `group:"flare"`
	}

	var flare Component
	app := fxtest.New(t,
		MockModule,
		fx.Provide(func() provides {
			return provides{
				Registration: FileRegistration("sub/dir/test.txt", func() (string, error) {
					return "hello, world", nil
				}),
			}
		}),
		fx.Populate(&flare),
	)

	defer app.RequireStart().RequireStop()

	content, err := flare.(Mock).GetFlareFile(t, "sub/dir/test.txt")
	require.NoError(t, err)
	require.Equal(t, "hello, world", content)
}
