// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package subscriptions

// subscriptionPoint implements SubscriptionPoint.
type SubscriptionPoint[M Message] struct {
	subscriptions []*subscription[M]
}

type subscription[M Message] struct {
	ch chan M
}

// NewSubscriptionPoint creates a new SubscriptionPoint.
//
// The given Subscriptions can be nil, in which case they are ignored.  This
// occurs when components that _might_ make a subscription choose not to (such
// as when those components are not enabled).
func NewSubscriptionPoint[M Message](subscriptions []Subscription[M]) *SubscriptionPoint[M] {
	// filter out the nil subscriptions, and cast the remainder to the concrete
	// type.
	concreteSubs := make([]*subscription[M], 0, len(subscriptions))
	for _, sub := range subscriptions {
		if sub != nil {
			concreteSubs = append(concreteSubs, sub.(*subscription[M]))
		}
	}
	return &SubscriptionPoint[M]{
		subscriptions: concreteSubs,
	}
}

// NewSubscription creates a new subscription of the required type.
func NewSubscription[M Message]() Subscription[M] {
	return &subscription[M]{make(chan M, 1)}
}

// Notify notifies all subscribers with a new message.
func (sp *SubscriptionPoint[M]) Notify(message M) {
	for _, s := range sp.subscriptions {
		s.send(message)
	}
}

// Chan implements Subscription#Chan.
func (s *subscription[M]) Chan() <-chan M {
	return s.ch
}

// send implements Subscription#send.
func (s *subscription[M]) send(message M) {
	s.ch <- message
}
