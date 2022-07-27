// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package agent implements a component representing the trace agent.
package agent

import (
	"fmt"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/processor"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/trace/api"
	"go.uber.org/fx"
)

type agent struct {
	processor processor.Component
}

type dependencies struct {
	fx.In

	Processor processor.Component
}

func newAgent(deps dependencies) Component {
	a := &agent{
		processor: deps.Processor,
	}

	// XXX temporary
	go func() {
		tkr := time.NewTicker(500 * time.Millisecond)
		for {
			<-tkr.C
			fmt.Printf("send\n")
			deps.Processor.PayloadChan() <- new(api.Payload)
		}
	}()

	return a
}
