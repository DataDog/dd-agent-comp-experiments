// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package ipcserver implements a component to manage the IPC API server and act
// as a client.
//
// The handlers in the HTTP server are supplied by other components, by providing a
// ipcserver.Route from their constructor.
//
// The Mock version of this component allows registration but does not actually
// start a server.
//
// XXX In a real agent, this would support TLS and gRPC and Auth and timeouts
// and stuff; see cmd/agent/api.
package ipcserver

import (
	"net/http"

	"go.uber.org/fx"
)

// team: agent-shared-components

const componentName = "comp/core/ipc/ipcserver"

// Component is the component type.
type Component interface {
}

// Mock implements mock-specific methods.
type Mock interface {
	Component

	// TODO: Get(path) ..
}

// Route is provided by other components in order to indicate routes that
// should be served via the IPC API.
type Route struct {
	fx.Out

	Route route `group:"ipcserver"`
}

// NewRoute creates a new Route instance for the named component.
func NewRoute(path string, handler http.HandlerFunc) Route {
	return Route{
		Route: route{path, handler},
	}
}

var Module = fx.Module(
	componentName,
	fx.Provide(newServer),
)

var MockModule = fx.Module(
	componentName,
	fx.Provide(newMock),
)
