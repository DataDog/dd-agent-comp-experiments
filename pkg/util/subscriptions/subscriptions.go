// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package subscriptions

// subscriptionPoint implements SubscriptionPoint.
type SubscriptionPoint[M Message] struct {
	subscribers []*subscription[M]
}

type subscription[M Message] struct {
	ch chan M
}

// NewSubscriptionPoint creates a new SubscriptionPoint.  The subscriptions
// must all be instances of `subscription`.
func NewSubscriptionPoint[M Message](subscriptions []Subscription[M]) *SubscriptionPoint[M] {
	subValues := make([]*subscription[M], len(subscriptions))
	for i, sub := range subscriptions {
		subValues[i] = sub.(*subscription[M])
	}
	return &SubscriptionPoint[M]{
		subscribers: subValues,
	}
}

// NewSubscription creates a new subscription of the required type.
func NewSubscription[M Message]() (Subscription[M], error) {
	return &subscription[M]{make(chan M, 1)}, nil
}

// Notify notifies all subscribers with a new message.
func (sp *SubscriptionPoint[M]) Notify(message M) {
	for _, s := range sp.subscribers {
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
