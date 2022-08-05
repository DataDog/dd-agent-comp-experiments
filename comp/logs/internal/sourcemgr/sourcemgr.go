// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package sourcemgr

import (
	"context"
	"sync"

	"github.com/djmitche/dd-agent-comp-experiments/comp/autodiscovery/scheduler"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/actor"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/subscriptions"
	"go.uber.org/fx"
)

type sourceMgr struct {
	// Mutex covers all fields in this type.
	sync.Mutex

	// started is true once the component has started
	started bool

	// sourceChangeTx is used to send SourceChanges
	sourceChangeTx subscriptions.Transmitter[SourceChange]

	// configChangeRx is used to subscribe to AD ConfigChanges
	configChangeRx subscriptions.Receiver[scheduler.ConfigChange]

	actor actor.Goroutine
}

type dependencies struct {
	fx.In

	Lc     fx.Lifecycle
	Params internal.BundleParams
	Config config.Component
	Pub    subscriptions.Publisher[SourceChange]
}

func newSourceMgr(deps dependencies) (Component, subscriptions.Subscription[scheduler.ConfigChange]) {
	sm := &sourceMgr{
		sourceChangeTx: deps.Pub.Transmitter(),
	}
	var sub subscriptions.Subscription[scheduler.ConfigChange]
	if deps.Params.ShouldStart(deps.Config) {
		sm.actor.HookLifecycle(deps.Lc, sm.run)
		deps.Lc.Append(fx.Hook{OnStart: sm.start})
		sub = subscriptions.NewSubscription[scheduler.ConfigChange]()
		sm.configChangeRx = sub.Receiver
	}
	return sm, sub
}

// start marks the component as started.
func (sm *sourceMgr) start(ctx context.Context) error {
	sm.Lock()
	defer sm.Unlock()
	sm.started = true
	return nil
}

func (sm *sourceMgr) run(ctx context.Context) {
	sources := map[string]*LogSource{}
	for {
		select {
		case chg := <-sm.configChangeRx.Chan():
			// XXX this temporarily subscribes to AD; there should be a scheduler in
			// between the two components.
			if chg.IsScheduled {
				src := &LogSource{Name: "logs-for-" + chg.Config.Name}
				sources[src.Name] = src
				sm.AddSource(src)
			} else {
				name := "logs-for-" + chg.Config.Name
				src := sources[name]
				if src != nil {
					sm.RemoveSource(src)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// AddSource implements Component#AddSource.
func (sm *sourceMgr) AddSource(source *LogSource) {
	sm.Lock()
	defer sm.Unlock()
	if !sm.started {
		panic("sourcemgr component has not been started")
	}
	sm.sourceChangeTx.Notify(SourceChange{
		IsAdd:  true,
		Source: source,
	})
}

// RemoveSource implements Component#RemoveSource.
func (sm *sourceMgr) RemoveSource(source *LogSource) {
	sm.Lock()
	defer sm.Unlock()
	if !sm.started {
		panic("sourcemgr component has not been started")
	}
	sm.sourceChangeTx.Notify(SourceChange{
		IsAdd:  false,
		Source: source,
	})
}
