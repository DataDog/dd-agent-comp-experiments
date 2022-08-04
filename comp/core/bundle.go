// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package core implements the "core" bundle, providing services common to all
// agent flavors and binaries.
//
// The constituent components serve as utilities and are mostly independent of
// one another.  Other components should depend on any components they need.
//
// This bundle does not depend on any other bundles.
package core

import (
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/flare"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/internal"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/ipc/ipcclient"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/ipc/ipcserver"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/log"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/status"
	"go.uber.org/fx"
)

// team: agent-shared-components

const componentName = "comp/core"

type BundleParams = internal.BundleParams

// Bundle defines the fx options for this component.
var Bundle = fx.Module(
	componentName,

	// apply defaults to BundleParams
	fx.Decorate(func(params *BundleParams) *BundleParams {
		if params != nil {
			return params
		}

		return &BundleParams{
			ConfFilePath: "", // config component applies its own defaults
			AutoStart:    false,
		}
	}),

	config.Module,
	flare.Module,
	health.Module,
	ipcclient.Module,
	ipcserver.Module,
	log.Module,
	status.Module,

	// instantiate the ipcserver unconditionally, as nothing else actually depends
	// on it (but it depends on a number of other things, such as flare and status)
	fx.Invoke(func(ipcserver.Component) {}),
)
