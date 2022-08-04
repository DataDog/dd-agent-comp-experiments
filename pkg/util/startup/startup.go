// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package startup

// AutoStart defines whether a component or bundle should start automatically.
type AutoStart int

const (
	// Never means that the component should never start.  This is the zero value.
	Never AutoStart = 0

	// Always means that the component should always start.
	Always AutoStart = 1

	// IfConfigured means that the component should consult its configuration
	// (e.g., `foo-agent.enabled`) to decide whether to start.  Components which
	// have no configuration treat this as Always
	IfConfigured = 2
)

// ShouldStart helps a component determine whether it should start; "enabled"
// is the component's configuration value, or `true` for components that are
// not controlled by a configuration parameter.
func (a AutoStart) ShouldStart(enabled bool) bool {
	switch a {
	case Always:
		return true
	case Never:
		return false
	case IfConfigured:
		return enabled
	default:
		return false // zero value is "Never", so use that as default
	}
}
