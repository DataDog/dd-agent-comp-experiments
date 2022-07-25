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
	"github.com/djmitche/dd-agent-comp-experiments/pkg/util/subscriptions"
	"go.uber.org/fx"
)

type launcher struct {
	log          log.Component
	subscription subscriptions.Subscriber[sourcemgr.SourceChange]
}

func (l *launcher) start(ctx context.Context) error {
	l.log.Debug("starting file launcher")
	go l.run()
	return nil
}

func (l *launcher) run() {
	// TODO: stop!
	for {
		select {
		case chg := <-l.subscription.Chan():
			l.log.Debug("got change", chg)
		}
	}
}

func newLauncher(lc fx.Lifecycle, log log.Component, sourcemgr sourcemgr.Component, mgr launchermgr.Component) (Component, error) {
	subscription, err := sourcemgr.Subscribe()
	if err != nil {
		return nil, err
	}
	l := &launcher{log, subscription}
	mgr.RegisterLauncher("file", l)
	lc.Append(fx.Hook{OnStart: l.start})
	return l, nil
}
