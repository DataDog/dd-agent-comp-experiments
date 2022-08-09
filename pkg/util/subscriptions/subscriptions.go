// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package subscriptions provides support for managing subscriptions between components.
//
// This package provides a simple interface with its Transmitter and Receiver types.
// Create Receivers with NewReceiver, and build a Transmitter to transmit to them.  Then
// send messages with tx.Notify() and receive them with <-rx.Chan().
//
// See the conventions documentation for a description of the component interface.
//
// Warning
//
// This package is not intended for high-bandwidth messaging such as metric
// samples.  It use should be limited to events that occur on a per-minute
// scale.
package subscriptions

import "go.uber.org/fx"

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

// Receiver defines a point where messages can be received
//
// A zero-valued receiver is valid, but will not receive messages.
type Receiver[M Message] struct {
	ch chan M
}

// NewReceiver creates a new Receiver.  Component-based subscriptions typically
// use NewSubscriber, instead.
func NewReceiver[M Message]() Receiver[M] {
	return Receiver[M]{
		ch: make(chan M, 1),
	}
}

// Chan gets the channel from which messages for this subscription should be read
func (s Receiver[M]) Chan() <-chan M {
	return s.ch
}

// Transmitter defines a point where messages can be sent
type Transmitter[M Message] struct {
	chs []chan M
}

// NewTransmitter creates a new Transmitter.  Component-based subscriptions
// typically use NewPublisher, instead.
//
// This ignores any zero-valued receivers.
func NewTransmitter[M Message](receivers []Receiver[M]) Transmitter[M] {
	// get the receivers' channels, filtering out nils
	chs := make([]chan M, 0, len(receivers))
	for _, rx := range receivers {
		if rx.ch != nil {
			chs = append(chs, rx.ch)
		}
	}
	return Transmitter[M]{chs}
}

// Notify notifies all associated receivers of a new message.
func (sp Transmitter[M]) Notify(message M) {
	for _, ch := range sp.chs {
		ch <- message
	}
}

// Subscription represents a component's request for a receiver of this type.
type Subscription[M Message] struct {
	fx.Out

	Receiver Receiver[M] `group:"subscriptions"`
}

// NewSubscription creates a new subscription of the required type.
//
// A receiving component's constructor should call this function, capture the
// Receiver field for later use, and return the Subscription.
func NewSubscription[M Message]() Subscription[M] {
	return Subscription[M]{
		Receiver: NewReceiver[M](),
	}
}

// Publisher represents a component's request for a transmitter of this type.
//
// A component's constructor should take an object of this type as an argument
// (or indirectly via an fx.In struct), and then call its Transmitter method.
type Publisher[M Message] struct {
	fx.In

	Receivers []Receiver[M] `group:"subscriptions"`
}

// Transmitter creates a transmitter for a publisher.
func (p Publisher[M]) Transmitter() Transmitter[M] {
	return NewTransmitter(p.Receivers)
}
