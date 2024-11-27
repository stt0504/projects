package pubsub

import (
	"context"
	"errors"
	"sync"
	"time"
)

type MyPubSub struct {
	mu            sync.Mutex
	topics        map[string][]chan interface{}
	subscriptions map[string][]*MySubscription
	closed        bool
	wg            sync.WaitGroup
}

type MySubscription struct {
	subj      string
	msgCh     chan interface{}
	pubSub    *MyPubSub
	cancelCtx context.Context
	cancelFn  context.CancelFunc
	wg        sync.WaitGroup
}

func NewPubSub() PubSub {
	return &MyPubSub{
		topics:        make(map[string][]chan interface{}),
		subscriptions: make(map[string][]*MySubscription),
	}
}

func (p *MyPubSub) Subscribe(subj string, cb MsgHandler) (Subscription, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, errors.New("pubsub is closed")
	}

	msgCh := make(chan interface{}, 100)

	ctx, cancelFn := context.WithCancel(context.Background())

	p.topics[subj] = append(p.topics[subj], msgCh)

	subscription := &MySubscription{
		subj:      subj,
		msgCh:     msgCh,
		pubSub:    p,
		cancelCtx: ctx,
		cancelFn:  cancelFn,
	}

	p.subscriptions[subj] = append(p.subscriptions[subj], subscription)

	go p.listenForMessages(subscription, cb)

	return subscription, nil
}

func (p *MyPubSub) listenForMessages(subscription *MySubscription, cb MsgHandler) {
	for {
		msg, ok := <-subscription.msgCh
		if !ok || subscription.cancelCtx.Err() != nil {
			return
		}
		cb(msg)
		p.wg.Done()
	}
}

func (p *MyPubSub) Publish(subj string, msg interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return errors.New("pubsub is closed")
	}

	subscribers, ok := p.topics[subj]
	if !ok || len(subscribers) == 0 {
		return errors.New("no subscribers for subject")
	}

	p.wg.Add(len(subscribers))

	for _, subscriberCh := range subscribers {
		select {
		case subscriberCh <- msg:
		default:
		}
	}

	return nil
}

func (s *MySubscription) Unsubscribe() {
	s.cancelFn()

	s.pubSub.mu.Lock()
	defer s.pubSub.mu.Unlock()

	subs := s.pubSub.subscriptions[s.subj]
	for i, sub := range subs {
		if sub == s {
			s.pubSub.subscriptions[s.subj] = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	subscribers := s.pubSub.topics[s.subj]
	for i, subscriberCh := range subscribers {
		if subscriberCh == s.msgCh {
			s.pubSub.topics[s.subj] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}

	close(s.msgCh)
}

func (p *MyPubSub) Close(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	time.Sleep(200 * time.Millisecond)

	select {
	case <-ctx.Done():
		for _, subs := range p.subscriptions {
			for _, sub := range subs {
				sub.cancelFn()
			}
		}
		return ctx.Err()
	default:
	}

	p.wg.Wait()
	p.closed = true

	for _, subs := range p.subscriptions {
		for _, sub := range subs {
			sub.cancelFn()
		}
	}

	for _, subscribers := range p.topics {
		for _, subscriberCh := range subscribers {
			close(subscriberCh)
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		p.wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
