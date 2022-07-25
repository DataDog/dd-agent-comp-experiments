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
	sp := NewSubscriptionPoint[string]()

	sub1, err := sp.Subscribe()
	require.NoError(t, err)
	sub2, err := sp.Subscribe()
	require.NoError(t, err)

	sp.Notify("hello!")
	require.Equal(t, "hello!", <-sub1.Chan())
	require.Equal(t, "hello!", <-sub2.Chan())

	sp.Unsubscribe(sub1)

	sp.Notify("goodbye!")
	require.Equal(t, 0, len(sub1.Chan()))
	require.Equal(t, "goodbye!", <-sub2.Chan())
}
