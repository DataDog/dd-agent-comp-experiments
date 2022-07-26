// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package ipcapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

type server struct {
	// disabled indicates that the component should do nothing.
	disabled bool

	// port is the port on which the server should run
	port int

	// router is the router used by the server
	router *mux.Router

	// server is the running server, if started and not disabled
	server *http.Server
}

type dependencies struct {
	fx.In
	Lc     fx.Lifecycle
	Params ModuleParams `optional:"true"`
	Config config.Component
}

func newServer(deps dependencies) Component {
	s := &server{
		disabled: deps.Params.Disabled,
		port:     deps.Config.GetInt("cmd_port"),
		router:   mux.NewRouter(),
	}
	deps.Lc.Append(fx.Hook{OnStart: s.start, OnStop: s.stop})
	return s
}

// Register implements Component#Register.
func (s *server) Register(path string, handler http.HandlerFunc) {
	if s.server != nil {
		panic("ipcapi component has already started")
	}

	if s.disabled {
		return
	}

	s.router.HandleFunc(path, handler)
}

// start starts the http server, if not disabled.
func (s *server) start(ctx context.Context) error {
	if s.disabled {
		return nil
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", s.port),
		Handler: s.router,
	}
	go s.server.ListenAndServe()
	return nil
}

// stop stops the http server, if not disabled.
func (s *server) stop(ctx context.Context) error {
	if s.disabled {
		return nil
	}

	if s.server != nil {
		defer func() { s.server = nil }()
		err := s.server.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
