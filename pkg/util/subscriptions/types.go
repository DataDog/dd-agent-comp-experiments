// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package subscriptions provides support for managing subscriptions between components.
//
// A subscription is a request from one component to another to be notified of events.
// Subscriptions are generic over the type of message being delivered.
//
// This package is not intended for high-bandwidth messaging.
package subscriptions

/* XXX Future Improvements
 *
 * This package is really simple as-is, but there are lots of additional
 * concerns to address:
 *  - Handling concurrent access (Subscribe, Unsubscribe, and Notify race in this version)
 *  - Handling consumers that fail to read from their channel
 *  - "Locking" subscriptions, as most components only allow changes to descriptions before
 *    they start.
 */

// Message is the type of the message handled by a subscription point.  It can be any type.
type Message interface{}

// Subscription defines a single subscriber to a SubscriptionPoint.
type Subscription[M Message] interface {
	// Chan returns the channel from which the subscriber should read messages.
	Chan() <-chan M

	// Send sends a message to the subscriber
	send(message M)
}
