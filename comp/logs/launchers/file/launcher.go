// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package file

import (
	"context"
	"time"

	"github.com/djmitche/dd-agent-comp-experiments/comp/core/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/health"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/log"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/internal/sourcemgr"
	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/launchermgr"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/actor"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/subscriptions"
	"go.uber.org/fx"
)

type launcher struct {
	log            log.Component
	sourceChangeRx subscriptions.Receiver[sourcemgr.SourceChange]
	actor          actor.Goroutine
	health         *health.Handle
}

type dependencies struct {
	fx.In

	Lc     fx.Lifecycle
	Config config.Component
	Params internal.BundleParams
	Log    log.Component
}

type provides struct {
	fx.Out

	Component
	HealthReg       health.Registration
	LauncherReg     launchermgr.Registration
	SourceChangeSub subscriptions.Subscription[sourcemgr.SourceChange]
}

func newLauncher(deps dependencies) provides {
	healthReg := health.NewRegistration(componentName)
	l := &launcher{
		log:    deps.Log,
		health: healthReg.Handle,
	}
	var sub subscriptions.Subscription[sourcemgr.SourceChange]
	if deps.Params.ShouldStart(deps.Config) {
		l.actor.HookLifecycle(deps.Lc, l.run)
		sub = subscriptions.NewSubscription[sourcemgr.SourceChange]()
		l.sourceChangeRx = sub.Receiver
	}
	return provides{
		Component:       l,
		HealthReg:       healthReg,
		LauncherReg:     launchermgr.NewRegistration("file", l),
		SourceChangeSub: sub,
	}
}

func (l *launcher) run(ctx context.Context) {
	monitor, stopMonitor := l.health.LivenessMonitor(time.Second)
	for {
		select {
		case chg := <-l.sourceChangeRx.Chan():
			l.log.Debug("file launcher got LogSource change", chg)
			// XXX start a tailer, etc. etc.
		case <-monitor:
		case <-ctx.Done():
			stopMonitor()
			return
		}
	}
}
