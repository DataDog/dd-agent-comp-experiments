// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package api

// Payload specifies information about a set of traces received by the API.
//
// XXX this would be the same as pkg/trace/api.Payload in the datadog-agent repo.
type Payload struct {
	Spans []Span
}

// Span is just a placeholder.
type Span struct {
	Data string
}
