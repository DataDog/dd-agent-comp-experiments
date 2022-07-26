// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package common

import "go.uber.org/fx"

// OneShot is an fx.Option that calls the given function, using fx.Invoke to
// load its arguments, and causes the app to shut down when the function
// returns.  It's used for one-shot agent subcommands like `agent status`.
func OneShot(f interface{}) fx.Option {
	return fx.Options(
		fx.Invoke(f),
		fx.Invoke(func(sd fx.Shutdowner) { sd.Shutdown() }),
	)
}
