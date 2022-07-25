// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package file

import (
	"context"

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
}

func (l *launcher) start(ctx context.Context) error {
	l.log.Debug("Starting file launcher")
	l.actor.Start(l.run)
	return nil
}

func (l *launcher) stop(ctx context.Context) error {
	l.log.Debug("Stopping file launcher")
	l.actor.Stop(context.Background())
	return nil
}

func (l *launcher) run(ctx context.Context) {
	for {
		select {
		case chg := <-l.subscription.Chan():
			l.log.Debug("got change", chg)
			// XXX start a tailer, etc. etc.
		case <-ctx.Done():
			return
		}
	}
}

func newLauncher(lc fx.Lifecycle, log log.Component, sourcemgr sourcemgr.Component, mgr launchermgr.Component) (Component, error) {
	subscription, err := sourcemgr.Subscribe()
	if err != nil {
		return nil, err
	}
	l := &launcher{
		log:          log,
		subscription: subscription,
	}
	mgr.RegisterLauncher("file", l)
	l.actor.HookLifecycle(lc, l.run)
	return l, nil
}
