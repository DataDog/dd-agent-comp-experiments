// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package flare

import (
	"testing"

	"github.com/djmitche/dd-agent-comp-experiments/comp/ipcapi"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFlareMechanics(t *testing.T) {
	var flare Component
	app := fxtest.New(t,
		Module,
		ipcapi.MockModule,
		fx.Populate(&flare),
	)

	flare.RegisterFile("test.txt", func() (string, error) {
		return "hello, world", nil
	})

	defer app.RequireStart().RequireStop()

	archiveFile, err := flare.CreateFlare()
	require.NoError(t, err)

	// XXX unzip archive file and verify..
	require.NotEqual(t, "", archiveFile)
}

func TestMock(t *testing.T) {
	var flare Component
	app := fxtest.New(t,
		MockModule,
		ipcapi.MockModule,
		fx.Populate(&flare),
	)

	flare.RegisterFile("sub/dir/test.txt", func() (string, error) {
		return "hello, world", nil
	})

	defer app.RequireStart().RequireStop()

	content, err := flare.(Mock).GetFlareFile(t, "sub/dir/test.txt")
	require.NoError(t, err)
	require.Equal(t, "hello, world", content)
}
