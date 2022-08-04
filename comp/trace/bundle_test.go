// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package trace

import (
	"testing"

	"github.com/djmitche/dd-agent-comp-experiments/comp/core"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestBundleDependencies(t *testing.T) {
	require.NoError(t, fx.ValidateApp(
		fx.Supply(core.BundleParams{}),
		core.Bundle,
		fx.Supply(BundleParams{}),
		Bundle))
}
