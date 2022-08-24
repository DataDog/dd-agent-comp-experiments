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
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/flare"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/health"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/internal"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcclient"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/ipc/ipcserver"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/log"
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/status"
	"go.uber.org/fx"
)

// team: agent-shared-components

const componentName = "comp/core"

type BundleParams = internal.BundleParams

// Bundle defines the fx options for this bundle.
var Bundle = fx.Module(
	componentName,

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

// MockBundle defines the mock fx options for this bundle.
var MockBundle = fx.Module(
	componentName,

	fx.Supply(internal.BundleParams{}),

	config.MockModule,
	flare.MockModule,
	health.Module,
	ipcclient.Module,
	ipcserver.MockModule,
	log.MockModule,
	status.Module,
)
