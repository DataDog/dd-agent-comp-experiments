// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package sourcemgr

import (
	"context"
	"sync"

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
}

func newSourceMgr(lc fx.Lifecycle) Component {
	sm := &sourceMgr{
		subscriptions: subscriptions.NewSubscriptionPoint[SourceChange](),
	}
	lc.Append(fx.Hook{OnStart: sm.start})
	return sm
}

// start marks the component as started.
func (sm *sourceMgr) start(ctx context.Context) error {
	sm.Lock()
	defer sm.Unlock()
	sm.started = true
	return nil
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
