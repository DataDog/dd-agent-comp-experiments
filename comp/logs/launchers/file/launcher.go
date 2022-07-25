// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package file

import (
	"context"

	"github.com/djmitche/dd-agent-comp-experiments/comp/logs/launchers/manager"
	"github.com/djmitche/dd-agent-comp-experiments/comp/util/log"
	"go.uber.org/fx"
)

type launcher struct {
	log log.Component
}

func (l *launcher) start(ctx context.Context) error {
	l.log.Debug("starting file launcher")
	return nil
}

func newLauncher(lc fx.Lifecycle, log log.Component, mgr manager.Component) Component {
	l := &launcher{log}
	mgr.RegisterLauncher("file", l)
	lc.Append(fx.Hook{OnStart: l.start})
	return l
}
