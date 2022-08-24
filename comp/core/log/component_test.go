// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package log

import (
	"testing"

	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/internal"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/comptest"
	"go.uber.org/fx"
)

func TestLogging(t *testing.T) {
	var log Component
	comptest.FxTest(t,
		fx.Supply(internal.BundleParams{}),
		config.MockModule,
		Module,
		fx.Populate(&log),
	).WithRunningApp(func() {
		log.Debug("hello, world.")
		// TODO: assert that the log succeeded
	})
}
