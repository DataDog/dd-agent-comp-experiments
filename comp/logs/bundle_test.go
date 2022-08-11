// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package logs

import (
	"testing"

	"github.com/DataDog/dd-agent-comp-experiments/comp/autodiscovery"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/launchers/file"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestBundleDependencies(t *testing.T) {
	require.NoError(t, fx.ValidateApp(
		fx.Supply(core.BundleParams{}),
		core.Bundle,
		fx.Supply(autodiscovery.BundleParams{}),
		autodiscovery.Bundle,

		// provide and require the launchers, since they are not required automatically
		file.Module,
		fx.Invoke(func(file.Component) {}),

		fx.Supply(BundleParams{}),
		Bundle))
}
