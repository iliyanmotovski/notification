package orm

import (
	"context"

	"github.com/latolukasz/beeorm"
)

type EventBroker interface {
	GetEventsConsumer(group string) EventsConsumer
}

type beeORMEventBroker struct {
	beeorm.EventBroker
}

func (b *beeORMEventBroker) GetEventsConsumer(group string) EventsConsumer {
	return &beeORMEventsConsumer{b.EventBroker.Consumer(group)}
}

type EventsConsumer interface {
	Consume(ctx context.Context, consumerIndex, prefetchCount int, handlerFunc interface{}) bool
}

type beeORMEventsConsumer struct {
	beeorm.EventsConsumer
}

func (b *beeORMEventsConsumer) Consume(ctx context.Context, consumerIndex, prefetchCount int, handlerFunc interface{}) bool {
	return b.EventsConsumer.ConsumeMany(ctx, consumerIndex, prefetchCount, handlerFunc.(func(events []beeorm.Event)))
}
