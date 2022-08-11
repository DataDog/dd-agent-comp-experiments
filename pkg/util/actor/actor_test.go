// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package actor

import (
	"context"
	"fmt"
	"time"
)

func Example() {
	type component struct {
		ch    chan int
		actor Actor
	}

	c := &component{
		ch: make(chan int),
	}

	run := func(ctx context.Context) {
		for {
			select {
			case v := <-c.ch:
				fmt.Printf("GOT: %d\n", v)
			case <-ctx.Done():
				fmt.Println("Stopping")
				return
			}
		}
	}

	c.actor.Start(run)
	c.ch <- 1
	c.ch <- 2
	c.actor.Stop(context.Background())

	// Output:
	// GOT: 1
	// GOT: 2
	// Stopping
}

func run(ctx context.Context) {
	tkr := time.NewTicker(time.Millisecond)
	for {
		select {
		case <-tkr.C:
			fmt.Println("tick")
		case <-ctx.Done():
			return
		}
	}
}
