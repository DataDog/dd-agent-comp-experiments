// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package autodiscovery implements discovery of configuration for dynamic
// entities like pods and containers.  It broadcasts changes to that configuration
// to subscribers.
package autodiscovery

import (
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/subscriptions"
	"go.uber.org/fx"
)

// team: container-integrations

const componentName = "comp/autodiscovery"

// Component is the component type.
type Component interface {
	// Subscribe registers a subscriber for schedule/unschedule events.  This
	// must be called before the component starts.  Subscribers cannot be
	// unsubscribed.
	Subscribe() (subscriptions.Subscriber[ConfigChange], error)
}

// Config defines config for a container or pod. XXX this is an
// integration.Config
type Config struct {
	Name string
}

// ConfigChange indicates a change to a config: being scheduled or unscheduled.
type ConfigChange struct {
	// IsScheduled is true when the configuration is scheduled, and false when it
	// is unscheduled.
	IsScheduled bool

	// Config is the config being changed.
	Config *Config
}

// Module defines the fx options for this component.
var Module fx.Option = fx.Module(
	componentName,
	fx.Provide(newAD),
)
