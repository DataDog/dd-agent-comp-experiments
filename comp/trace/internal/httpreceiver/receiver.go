// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package httpreceiver

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	"github.com/djmitche/dd-agent-comp-experiments/comp/config"
	"github.com/djmitche/dd-agent-comp-experiments/comp/trace/internal/processor"
	"github.com/djmitche/dd-agent-comp-experiments/pkg/trace/api"
	"go.uber.org/fx"
)

type receiver struct {
	// port is the port on which the server is running.
	port int

	// server is the running server
	server *http.Server

	// processorChan is the channel to the processor component
	processorChan chan<- *api.Payload
}

type dependencies struct {
	fx.In

	Lc        fx.Lifecycle
	Config    config.Component
	Processor processor.Component
}

func newReceiver(deps dependencies) Component {
	r := &receiver{
		port:          deps.Config.GetInt("apm_config.receiver_port"),
		processorChan: deps.Processor.PayloadChan(),
	}
	deps.Lc.Append(fx.Hook{OnStart: r.start, OnStop: r.stop})
	return r
}

// start starts the http server.
func (r *receiver) start(ctx context.Context) error {
	r.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", r.port),
		Handler: http.HandlerFunc(r.handler),
	}
	go r.server.ListenAndServe()
	return nil
}

// handle handles an HTTP request.
func (r *receiver) handler(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		return
	}

	// XXX: this is not AT ALL what spans look like, but it's easy to use with
	// curl for demo purposes.

	scanner := bufio.NewScanner(req.Body)
	spans := []api.Span{}
	for scanner.Scan() {
		spans = append(spans, api.Span{Data: scanner.Text()})
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("uhoh: %s\n", err)
		w.WriteHeader(400)
	}

	r.processorChan <- &api.Payload{Spans: spans}
}

// stop stops the http server.
func (r *receiver) stop(ctx context.Context) error {
	if r.server != nil {
		defer func() { r.server = nil }()
		err := r.server.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
