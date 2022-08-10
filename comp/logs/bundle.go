// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package logs collects the packages related to the logs agent.
//
// The comp/logs.Bundle options include most components in this bundle, but
// does not include launchers.  Applications must include the desired launchers
// itself.  This is done to support smaller binary sizes for applications that
// do not require all launchers.
//
// This bundle depends on comp/core and comp/autodiscovery.
package logs

import (
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/agent"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/internal/sourcemgr"
	"github.com/DataDog/dd-agent-comp-experiments/comp/logs/launchers/launchermgr"
	"go.uber.org/fx"
)

// team: agent-metrics-logs

const componentName = "comp/logs"

// BundleParams defines the parameters for this bundle.
type BundleParams = internal.BundleParams

// Bundle defines the fx options for this bundle.
var Bundle fx.Option = fx.Module(
	componentName,

	agent.Module,
	launchermgr.Module,
	sourcemgr.Module,

	// always instantiate the controlling component; it will decide whether to
	// start based on AutoStart.
	fx.Invoke(func(agent.Component) {}),
)
