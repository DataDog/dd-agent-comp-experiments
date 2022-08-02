// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package sourcemgr

import (
	"context"
	"sync"

	"github.com/djmitche/dd-agent-comp-experiments/comp/autodiscovery"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/actor"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/subscriptions"
	"go.uber.org/fx"
)

type sourceMgr struct {
	// Mutex covers all fields in this type.
	sync.Mutex

	// started is true once the component has started
	started bool

	// subscriptions contains the subcriptions for source additions/removals
	subscriptions subscriptions.SubscriptionPoint[SourceChange]

	// subscription is the subscription to AD
	subscription subscriptions.Subscriber[autodiscovery.ConfigChange]

	actor actor.Goroutine
}

func newSourceMgr(lc fx.Lifecycle, autodiscovery autodiscovery.Component) (Component, error) {
	subscription, err := autodiscovery.Subscribe()
	if err != nil {
		return &sourceMgr{}, err
	}
	sm := &sourceMgr{
		subscriptions: subscriptions.NewSubscriptionPoint[SourceChange](),
		subscription:  subscription,
	}
	sm.actor.HookLifecycle(lc, sm.run)
	lc.Append(fx.Hook{OnStart: sm.start})
	return sm, nil
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
		case chg := <-sm.subscription.Chan():
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

// Subscribe implements Component#Subscribe.
func (sm *sourceMgr) Subscribe() (subscriptions.Subscriber[SourceChange], error) {
	sm.Lock()
	defer sm.Unlock()
	if sm.started {
		panic("sourcemgr component has already been started")
	}
	return sm.subscriptions.Subscribe()
}

// AddSource implements Component#AddSource.
func (sm *sourceMgr) AddSource(source *LogSource) {
	sm.Lock()
	defer sm.Unlock()
	if !sm.started {
		panic("sourcemgr component has not been started")
	}
	sm.subscriptions.Notify(SourceChange{
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
	sm.subscriptions.Notify(SourceChange{
		IsAdd:  false,
		Source: source,
	})
}
