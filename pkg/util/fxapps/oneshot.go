// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package fxapps

import (
	"context"
	"reflect"

	"go.uber.org/fx"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// OneShot runs the given function in an fx.App using the supplied options.
// The function's arguments are supplied by Fx and can be any provided type.
// The function must return `error.
//
// The resulting app starts all components, then invokes the function, then
// immediately shuts down.  This is typically used for command-line tools like
// `agent status`.
func OneShot(oneShotFunc interface{}, opts ...fx.Option) error {
	ftype := reflect.TypeOf(oneShotFunc)
	if ftype == nil || ftype.Kind() != reflect.Func {
		panic("OneShot requires a function as its first argument")
	}

	// verify it returns error
	if ftype.NumOut() != 1 || !ftype.Out(0).Implements(errorInterface) {
		panic("OneShot function must return error or nothing")
	}

	// build an function with the same signature as oneShotFunc that will
	// capture the args and do nothing.
	var oneShotArgs []reflect.Value
	captureArgs := reflect.MakeFunc(
		ftype,
		func(args []reflect.Value) []reflect.Value {
			oneShotArgs = args
			// return a single nil value of type error
			return []reflect.Value{reflect.Zero(errorInterface)}
		})
	// fx.Invoke that function to capture the args at startup
	opt := fx.Invoke(captureArgs.Interface())
	opts = append(opts, opt)
	app := fx.New(opts...)

	startCtx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		return err
	}

	// call the original oneShotFunc with the args captured during
	// app startup
	res := reflect.ValueOf(oneShotFunc).Call(oneShotArgs)
	if !res[0].IsNil() {
		err := res[0].Interface().(error)
		return err
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout())
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		return err
	}

	return nil
}
