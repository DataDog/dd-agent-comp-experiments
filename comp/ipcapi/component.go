// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package ipcapi implements a component to manage the IPC API server and act
// as a client.
//
// It allows other components to register handlers, and manages
// startup and shutdown of the HTTP server.
//
// It also supports simple GET requests to the server.
//
// The Mock version of this component allows registration but does not actually
// start a server, and does not require ModuleParams.
//
// XXX In a real agent, this would support TLS and gRPC and Auth and timeouts
// and stuff; see cmd/agent/api.
package ipcapi

import (
	"net/http"

	"go.uber.org/fx"
)

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// Register registers a handler at an HTTP path.
	Register(path string, handler http.HandlerFunc)

	// GetJSON gets the body of the response, as JSON
	GetJSON(path string, v any) error
}

// Mock implements mock-specific methods.
type Mock interface {
	Component

	// TODO: Get(path) ..
}

type ModuleParams struct {
	// Disabled indicates that the component should ignore all registration and
	// perform no monitoring.  This is intended for one-shot processes such as
	// `agent status`.
	Disabled bool
}

var Module = fx.Module(
	"comp/ipcapi",
	fx.Provide(newIpcAPI),
)

var MockModule = fx.Module(
	"comp/ipcapi",
	fx.Provide(newMock),
)
