// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package internal

// BundleParams must be defined here to avoid package dependency cycles.

// BundleParams defines the parameters for this bundle.
type BundleParams struct {
	// AutoStart determines whether AutoDiscovery should start automatically,
	// defaulting to false.
	AutoStart bool
}
