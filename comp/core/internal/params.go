// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package internal

import (
	"github.com/DataDog/dd-agent-comp-experiments/pkg/util/startup"
)

// BundleParams must be defined here to avoid package dependency cycles.

// BundleParams defines the parameters for this bundle.
type BundleParams struct {
	// ConfFilePath is the path to the configuration file.
	ConfFilePath string

	// AutoStart determines whether components in this bundle should start
	// automatically.  This is typically true for long-running processes and
	// false for one-shot processes.  This defaults to Always.
	AutoStart startup.AutoStart

	// Console determines whether log messages should be output to the console.
	Console bool
}

// ShouldStart determines whether the bundle should start.
func (p BundleParams) ShouldStart() bool {
	return p.AutoStart.ShouldStart(true)
}
