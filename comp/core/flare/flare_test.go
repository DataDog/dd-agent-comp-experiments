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

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/comptest"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/startup"
	"github.com/mholt/archiver"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestFlareMechanics(t *testing.T) {
	flareDir := t.TempDir()

	var flare Component
	comptest.FxTest(t,
		Module,
		log.Module,
		config.MockModule,
		fx.Supply(internal.BundleParams{AutoStart: startup.Always}),
		fx.Provide(func() Registration {
			return FileRegistration("greeting.txt", func() (string, error) {
				return "hello, world", nil
			})
		}),
		fx.Populate(&flare),
	).WithRunningApp(func() {
		archiveFile, err := flare.CreateFlare()
		require.NoError(t, err)
		// this will create a temporary file without t.TempDir, so we must clean it
		// up manually
		defer os.RemoveAll(archiveFile)

		require.NotEqual(t, "", archiveFile)
		err = archiver.Extract(archiveFile, "hostname/greeting.txt", flareDir)
		require.NoError(t, err)

		content, err := ioutil.ReadFile(filepath.Join(flareDir, "hostname", "greeting.txt"))
		require.NoError(t, err)
		require.Equal(t, "hello, world", string(content))
	})
}

func TestMock(t *testing.T) {
	var flare Component
	comptest.FxTest(t,
		MockModule,
		fx.Provide(func() Registration {
			return FileRegistration("sub/dir/test.txt", func() (string, error) {
				return "hello, world", nil
			})
		}),
		fx.Populate(&flare),
	).WithRunningApp(func() {
		content, err := flare.(Mock).GetFlareFile(t, "sub/dir/test.txt")
		require.NoError(t, err)
		require.Equal(t, "hello, world", content)
	})
}
