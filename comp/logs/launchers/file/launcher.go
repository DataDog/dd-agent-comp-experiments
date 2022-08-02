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
	subscription subscriptions.Subscription[sourcemgr.SourceChange]
	actor        actor.Goroutine
	health       *health.Registration
}

type dependencies struct {
	fx.In

	Lc  fx.Lifecycle
	Log log.Component
}

type provides struct {
	fx.Out

	Component
	HealthReg      *health.Registration      `group:"health"`
	Subscription   sourcemgr.Subscription    `group:"sourcemgr"`
	LauncherMgrReg *launchermgr.Registration `group:"launchermgr"`
}

func newLauncher(deps dependencies) (provides, error) {
	subscription, err := sourcemgr.Subscribe()
	if err != nil {
		return provides{}, err
	}
	l := &launcher{
		log:          deps.Log,
		subscription: subscription,
		health:       health.NewRegistration(componentName),
	}
	l.actor.HookLifecycle(deps.Lc, l.run)
	return provides{
		Component:      l,
		HealthReg:      l.health,
		LauncherMgrReg: launchermgr.NewRegistration("file", l),
		Subscription:   subscription,
	}, nil
}

func (l *launcher) run(ctx context.Context) {
	monitor, stopMonitor := l.health.LivenessMonitor(time.Second)
	for {
		select {
		case chg := <-l.subscription.Chan():
			l.log.Debug("file launcher got LogSource change", chg)
			// XXX start a tailer, etc. etc.
		case <-monitor:
		case <-ctx.Done():
			stopMonitor()
			return
		}
	}
}
