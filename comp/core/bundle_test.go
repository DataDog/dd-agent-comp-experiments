// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package core

import (
	"testing"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/flare"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/health"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcclient"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcserver"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/status"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestBundleDependencies(t *testing.T) {
	require.NoError(t, fx.ValidateApp(
		// instantiate all of the core components, since this is not done
		// automatically.
		fx.Invoke(func(config.Component) {}),
		fx.Invoke(func(flare.Component) {}),
		fx.Invoke(func(health.Component) {}),
		fx.Invoke(func(ipcclient.Component) {}),
		fx.Invoke(func(ipcserver.Component) {}),
		fx.Invoke(func(log.Component) {}),
		fx.Invoke(func(status.Component) {}),

		fx.Supply(BundleParams{}),
		Bundle))
}
