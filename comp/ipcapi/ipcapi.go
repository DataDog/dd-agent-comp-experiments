// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package ipcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

type ipcapi struct {
	// disabled indicates that the component should do nothing.
	disabled bool

	// port is the port on which the server is running.
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
	Routes []*Route
}

func newIpcAPI(deps dependencies) Component {
	a := &ipcapi{
		disabled: deps.Params.Disabled,
		port:     deps.Config.GetInt("cmd_port"),
		router:   mux.NewRouter(),
	}

	for _, r := range deps.Routes {
		a.router.HandleFunc(r.path, r.handler)
	}

	deps.Lc.Append(fx.Hook{OnStart: a.start, OnStop: a.stop})
	return a
}

func newMock() Component {
	a := &ipcapi{
		disabled: true,
		port:     0,
		router:   mux.NewRouter(),
	}
	return a
}

// GetJSON implements Component#GetJSON.
func (a *ipcapi) GetJSON(path string, v any) error {
	url := fmt.Sprintf("http://127.0.0.1:%d%s", a.port, path)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error contacting Agent: %s", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Error contacting Agent: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)

	if err != nil {
		return fmt.Errorf("Error decoding Agent response: %s", err)
	}

	return nil

}

// start starts the http server, if not disabled.
func (a *ipcapi) start(ctx context.Context) error {
	if a.disabled {
		return nil
	}

	a.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", a.port),
		Handler: a.router,
	}
	go a.server.ListenAndServe()
	return nil
}

// stop stops the http server, if not disabled.
func (a *ipcapi) stop(ctx context.Context) error {
	if a.disabled {
		return nil
	}

	if a.server != nil {
		defer func() { a.server = nil }()
		err := a.server.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
