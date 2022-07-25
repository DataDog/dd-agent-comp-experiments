// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package health

// SimpleRegistration is the result of registering with RegisterSimple.
type SimpleRegistration struct {
	health    *health
	component string
}

// SetUnhealthy records this component as being unhealthy, with the included message
// summarizing the problem.  This can be called at any time.
func (reg *SimpleRegistration) SetUnhealthy(message string) {
	reg.health.setHealth(reg.component, false, message)
}

// SetHealthy records this component as being healthy.  This can be called at any time.
func (reg *SimpleRegistration) SetHealthy() {
	reg.health.setHealth(reg.component, true, "")
}
