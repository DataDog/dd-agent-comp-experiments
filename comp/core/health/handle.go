// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

// Handle is the interface from other components to the health component.
//
// Handle methods must not be called until the calling component has
// started.
type Handle struct {
	// component is the name of the component being monitored.
	component string

	// health links to the comp/core/health component, once registration is
	// complete.
	health *health
}

// SetUnhealthy records this component as being unhealthy, with the included message
// summarizing the problem.
//
// This method must not be called before the monitored component has started.
func (reg *Handle) SetUnhealthy(message string) {
	// if comp/core/health hasn't been created, then there is nothing to do.
	if reg.health != nil {
		reg.health.setHealth(reg.component, false, message)
	}
}

// SetHealthy records this component as being healthy.
//
// This method must not be called before the monitored component has started.
func (reg *Handle) SetHealthy() {
	// if comp/core/health hasn't been created, then there is nothing to do.
	if reg.health != nil {
		reg.health.setHealth(reg.component, true, "")
	}
}
