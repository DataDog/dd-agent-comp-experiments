// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package log

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

type mock struct {
	sync.Mutex
	t         *testing.T
	capturing bool
	captured  []string
}

func newMockLogger(t *testing.T) Component {
	return &mock{
		t: t,
	}
}

// Configure implements Component#Configure.
func (m *mock) Configure(level string) error {
	return nil
}

// Debug implements Component#Debug.
func (m *mock) Debug(v ...interface{}) {
	m.Lock()
	defer m.Unlock()

	if m.t != nil {
		m.t.Log(v...)
	}

	if m.capturing {
		var bldr strings.Builder
		_, _ = fmt.Fprintln(&bldr, v...)
		m.captured = append(m.captured, bldr.String())
	}
}

// Flush implements Component#Flush.
func (*mock) Flush() {
}

// StartCapture implements Mock#StartCapture.
func (m *mock) StartCapture() {
	m.Lock()
	defer m.Unlock()

	m.capturing = true
	m.captured = nil
}

// Captured implements Mock#Captured.
func (m *mock) Captured() []string {
	m.Lock()
	defer m.Unlock()

	// return a copy of the captured log messages, to avoid concurrent access
	// to the slice
	return append([]string{}, m.captured...)
}

// EndCapture implements Mock#EndCapture.
func (m *mock) EndCapture() {
	m.Lock()
	defer m.Unlock()

	m.capturing = false
	m.captured = nil
}
