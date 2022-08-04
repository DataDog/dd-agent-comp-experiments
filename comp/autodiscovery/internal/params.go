// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package internal

import (
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/startup"
)

// BundleParams must be defined here to avoid package dependency cycles.

// BundleParams defines the parameters for this bundle.
type BundleParams struct {
	// AutoStart determines whether AutoDiscovery should start automatically,
	// defaulting to false.
	AutoStart startup.AutoStart
}

// ShouldStart determines whether the bundle should start.
func (p BundleParams) ShouldStart() bool {
	return p.AutoStart.ShouldStart(true)
}
