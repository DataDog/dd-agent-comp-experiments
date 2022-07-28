// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package file

import (
	"context"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal/sourcemgr"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/launchermgr"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/actor"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/subscriptions"
	"go.uber.org/fx"
)

type launcher struct {
	log          log.Component
	subscription subscriptions.Subscriber[sourcemgr.SourceChange]
	actor        actor.Goroutine
	health       *health.ActorRegistration
}

type dependencies struct {
	fx.In

	Lc          fx.Lifecycle
	Log         log.Component
	Sourcemgr   sourcemgr.Component
	LauncherMgr launchermgr.Component
	Health      health.Component
}

func newLauncher(deps dependencies) (Component, error) {
	subscription, err := deps.Sourcemgr.Subscribe()
	if err != nil {
		return nil, err
	}
	l := &launcher{
		log:          deps.Log,
		subscription: subscription,
		health:       deps.Health.RegisterActor("comp/logs/launchers/file", 1*time.Second),
	}
	deps.LauncherMgr.RegisterLauncher("file", l)
	l.actor.HookLifecycle(deps.Lc, l.run)
	return l, nil
}

func (l *launcher) run(ctx context.Context) {
	defer l.health.Stop()
	for {
		select {
		case chg := <-l.subscription.Chan():
			l.log.Debug("got change", chg)
			// XXX start a tailer, etc. etc.
		case <-l.health.Chan():
		case <-ctx.Done():
			return
		}
	}
}
