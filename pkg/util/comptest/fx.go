// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package comptest

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// IsTest is true if this is a test run.  This will always be false in
// "real" binaries.
var IsTest bool

// FxTester wraps fxtest.App with some convenience methods to help
// testing components.
type FxTester struct {
	*fxtest.App

	t testing.TB
}

// FxTest creates a new FxTester.  Its arguments match those of fxtest.New.

// Aside from the methods on FxTester, this adds the given testing.TB to the
// provided types and sets `IsTest` to true.
//
// Components can use IsTest to check that undesired functionality isn't
// accidentally enabled in tests.
func FxTest(t testing.TB, opts ...fx.Option) *FxTester {
	// set the IsTest flag for the duration of this test
	IsTest = true
	t.Cleanup(func() { IsTest = false })

	opts = append(opts, fx.Supply(t))
	return &FxTester{
		App: fxtest.New(t.(fxtest.TB),
			opts...),
		t: t,
	}
}

// WithRunningApp runs the given function after starting the app.
func (f *FxTester) WithRunningApp(test func()) {
	defer f.RequireStart().RequireStop()
	test()
}
