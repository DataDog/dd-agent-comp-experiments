// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package internal

import (
	"github.com/DataDog/dd-agent-comp-experiments/comp/core/config"
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/startup"
)

// BundleParams must be defined here to avoid package dependency cycles.

// BundleParams defines the parameters for this bundle.
type BundleParams struct {
	// AutoStart determines whether trace-agent components should start, defaulting
	// to Never.
	AutoStart startup.AutoStart
}

// ShouldStart determines whether the bundle should start, based on
// configuration.
func (p BundleParams) ShouldStart(config config.Component) bool {
	return p.AutoStart.ShouldStart(config.GetBool("apm_config.enabled"))
}
