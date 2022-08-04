// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package subscriptions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscribeUnsubscribe(t *testing.T) {
	sub1 := NewSubscription[string]()
	sub2 := NewSubscription[string]()
	sp := NewSubscriptionPoint[string]([]Subscription[string]{sub1, sub2})

	sp.Notify("hello!")
	require.Equal(t, "hello!", <-sub1.Chan())
	require.Equal(t, "hello!", <-sub2.Chan())

	require.Equal(t, 0, len(sub1.Chan()))
	require.Equal(t, 0, len(sub2.Chan()))
}
