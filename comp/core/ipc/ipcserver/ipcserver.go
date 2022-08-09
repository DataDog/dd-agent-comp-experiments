// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package ipcserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/djmitche/dd-agent-comp-experiments/comp/core/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/core/internal"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

type server struct {
	// autoStart indicates that the component should do nothing.
	autoStart bool

	// port is the port on which the server is running.
	port int

	// router is the router used by the server
	router *mux.Router

	// server is the running server, if started
	server *http.Server
}

// route is provided by other components in order to indicate routes that
// should be served via the IPC API.
type route struct {
	// path is the path at which this handler should be registered
	path string

	// handler is the handler for this path.
	handler http.HandlerFunc
}

type dependencies struct {
	fx.In
	Lc     fx.Lifecycle
	Params internal.BundleParams
	Config config.Component
	Routes []route `group:"ipcserver"`
}

func newServer(deps dependencies) Component {
	a := &server{
		autoStart: deps.Params.ShouldStart(),
		port:      deps.Config.GetInt("cmd_port"),
		router:    mux.NewRouter(),
	}

	for _, r := range deps.Routes {
		a.router.HandleFunc(r.path, r.handler)
	}

	deps.Lc.Append(fx.Hook{OnStart: a.start, OnStop: a.stop})
	return a
}

func newMock() Component {
	a := &server{
		autoStart: false,
		port:      0,
		router:    mux.NewRouter(),
	}
	return a
}

// start starts the http server, if autoStart is true.
func (a *server) start(ctx context.Context) error {
	if !a.autoStart {
		return nil
	}

	a.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", a.port),
		Handler: a.router,
	}
	go a.server.ListenAndServe()
	return nil
}

// stop stops the http server, if started
func (a *server) stop(ctx context.Context) error {
	if a.server != nil {
		defer func() { a.server = nil }()
		err := a.server.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
