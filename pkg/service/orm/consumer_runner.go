package orm

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/iliyanmotovski/notification/pkg/queue"
	"github.com/latolukasz/beeorm"
)

const (
	obtainLockRetryDuration = time.Second
)

type ConsumerHandler func(ormService Engine, entityIDs []uint64) error

type RunnerBeeORM struct {
	wg       sync.WaitGroup
	mu       sync.Mutex
	ctx      context.Context
	registry RegistryService
}

func NewConsumerRunner(ctx context.Context, registry RegistryService) *RunnerBeeORM {
	return &RunnerBeeORM{ctx: ctx, registry: registry}
}

func (r *RunnerBeeORM) RunConsumerMany(consumerHandler ConsumerHandler, queueName string, prefetchCount int) {
	r.wg.Add(1)

	go func(r *RunnerBeeORM, consumerHandler ConsumerHandler, queueName string, prefetchCount int) {
		consumerGroupName := queue.GetConsumerGroupName(queueName)

		ormService := r.registry.GetORMService()
		eventsConsumer := ormService.GetEventBroker().GetEventsConsumer(consumerGroupName)

		cacheService := ormService.GetCacheService(StreamsPool)

		r.mu.Lock()
		currentConsumerIndex := addConsumerGroup(cacheService, consumerGroupName)
		r.mu.Unlock()

		log.Printf("RunConsumerMany initialized (%s) index %d", queueName, currentConsumerIndex)

		for {
			// eventsConsumer.Consume should block and not return anything
			// if it returns true => this consumer is exited with no errors when we cancel the context
			// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry
			if exitedWithNoErrors := eventsConsumer.Consume(r.ctx, currentConsumerIndex, prefetchCount, func(events []beeorm.Event) {
				log.Printf("%d new dirty events in %s", len(events), queueName)

				entityIDs := make([]uint64, len(events))

				for i, event := range events {
					entityID := beeorm.EventDirtyEntity(event).ID()
					entityIDs[i] = entityID
				}

				if err := consumerHandler(ormService, entityIDs); err != nil {
					r.mu.Lock()
					removeConsumerGroup(cacheService, consumerGroupName, currentConsumerIndex)
					r.mu.Unlock()

					r.wg.Done()
					panic(err)
				}

				log.Printf("consumed %d dirty events in %s", len(events), queueName)
			}); !exitedWithNoErrors {
				log.Printf("RunConsumerMany failed to start (%s) - retrying in %.1f seconds", queueName, obtainLockRetryDuration.Seconds())
				time.Sleep(obtainLockRetryDuration)

				continue
			}

			log.Println("eventsConsumer.Consume returned true")
			log.Printf("RunConsumerMany exited (%s) index: %d", queueName, currentConsumerIndex)

			break
		}

		r.mu.Lock()
		removeConsumerGroup(cacheService, consumerGroupName, currentConsumerIndex)
		r.mu.Unlock()

		r.wg.Done()
	}(r, consumerHandler, queueName, prefetchCount)
}

func (r *RunnerBeeORM) Wait() {
	r.wg.Wait()
}

const consumerGroupsKey = "consumer_groups"

type indexer struct {
	LatestIndex           int
	ActiveConsumerIndexes map[int]struct{}
}

func addConsumerGroup(cacheService CacheService, consumerGroupName string) int {
	indexerValue, err := getConsumerGroupIndexer(cacheService, consumerGroupName)
	if err != nil {
		panic(err)
	}

	if indexerValue == nil {
		indexerValue = &indexer{}
	}

	if indexerValue.ActiveConsumerIndexes == nil {
		indexerValue.ActiveConsumerIndexes = map[int]struct{}{}
	}

	indexerValue.LatestIndex++
	indexerValue.ActiveConsumerIndexes[indexerValue.LatestIndex] = struct{}{}

	err = setConsumerGroupIndexer(cacheService, consumerGroupName, indexerValue)
	if err != nil {
		panic(err)
	}

	return indexerValue.LatestIndex
}

func removeConsumerGroup(cacheService CacheService, consumerGroupName string, indexToRemove int) {
	indexerValue, err := getConsumerGroupIndexer(cacheService, consumerGroupName)
	if err != nil {
		panic(err)
	}

	delete(indexerValue.ActiveConsumerIndexes, indexToRemove)

	err = setConsumerGroupIndexer(cacheService, consumerGroupName, indexerValue)
	if err != nil {
		panic(err)
	}
}

func setConsumerGroupIndexer(cacheService CacheService, consumerGroupName string, indexer *indexer) error {
	marshaled, err := json.Marshal(indexer)
	if err != nil {
		return err
	}

	cacheService.Set(consumerGroupsKey, consumerGroupName, marshaled)

	return err
}

func getConsumerGroupIndexer(cacheService CacheService, consumerGroupName string) (*indexer, error) {
	marshaled, has := cacheService.Get(consumerGroupsKey, consumerGroupName)
	if !has {
		return nil, nil
	}

	indexer := &indexer{}
	if err := json.Unmarshal([]byte(marshaled), indexer); err != nil {
		return nil, err
	}

	return indexer, nil
}
