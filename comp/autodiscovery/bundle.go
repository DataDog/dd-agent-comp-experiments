// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package autodiscovery implements autodiscovery: collecting integration configuration
// and monitoring for updates.
//
// This bundle depends on comp/core.
package autodiscovery

import (
	"github.com/djmitche/dd-agent-comp-experiments/comp/autodiscovery/internal"
	"github.com/djmitche/dd-agent-comp-experiments/comp/autodiscovery/scheduler"
	"go.uber.org/fx"
)

// team: container-integrations

const componentName = "comp/autodiscovery"

type BundleParams = internal.BundleParams

// Bundle defines the fx options for this bundle.
var Bundle = fx.Module(
	componentName,

	scheduler.Module,

	// instantiate the scheduler unconditionally, as nothing else actually depends
	// on it.
	fx.Invoke(func(scheduler.Component) {}),
)
