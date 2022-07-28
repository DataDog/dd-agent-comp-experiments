// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package agent implements a component representing the trace agent.
package agent

import (
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/httpreceiver"
	"go.uber.org/fx"
)

type agent struct {
}

type dependencies struct {
	fx.In

	HTTPReceiver httpreceiver.Component // required just to load the component
}

func newAgent(deps dependencies) Component {
	a := &agent{}

	// TODO: this will likely carry a reference to Receiver, Processor, and so
	// on to handle requests for Status, stats, etc.

	return a
}
