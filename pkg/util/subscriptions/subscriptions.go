// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package subscriptions

// subscriptionPoint implements SubscriptionPoint.
type subscriptionPoint[M Message] struct {
	subscribers []*subscriber[M]
}

type subscriber[M Message] struct {
	ch chan M
}

// NewSubscriptionPoint creates a new SubscriptionPoint.
func NewSubscriptionPoint[M Message]() SubscriptionPoint[M] {
	return &subscriptionPoint[M]{}
}

// Subscribe implements SubscriptionPoint#Subscribe.
func (sp *subscriptionPoint[M]) Subscribe() (Subscriber[M], error) {
	sub := &subscriber[M]{make(chan M, 1)}
	sp.subscribers = append(sp.subscribers, sub)
	return sub, nil
}

// Unsubscribe implements SubscriptionPoint#Unsubscribe.
func (sp *subscriptionPoint[M]) Unsubscribe(sub Subscriber[M]) {
	for i, s := range sp.subscribers {
		if s == sub {
			sp.subscribers = append(sp.subscribers[:i], sp.subscribers[i+1:]...)
			return
		}
	}
}

// Notify implements SubscriptionPoint#Notify.
func (sp *subscriptionPoint[M]) Notify(message M) {
	for _, s := range sp.subscribers {
		s.send(message)
	}
}

// Chan implements Subscriber#Chan.
func (s *subscriber[M]) Chan() <-chan M {
	return s.ch
}

// send implements Subscriber#send.
func (s *subscriber[M]) send(message M) {
	s.ch <- message
}
