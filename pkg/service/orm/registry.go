package orm

import (
	"errors"
	"fmt"
	"github.com/iliyanm/notification/pkg/queue"

	entitybeeorm "github.com/iliyanm/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanm/notification/pkg/service/config"
	"github.com/latolukasz/beeorm"
)

const (
	streamsPool = "streams_pool"
)

type RegistryService interface {
	GetORMService() Engine
}

func NewORMRegistryService(configService config.Config) (RegistryService, func(), error) {
	registry := beeorm.NewRegistry()

	configuration, ok := configService.Get("orm")
	if !ok {
		return nil, nil, errors.New("no orm config")
	}

	yamlConfig := map[string]interface{}{}
	for k, v := range configuration.(map[interface{}]interface{}) {
		yamlConfig[fmt.Sprint(k)] = v
	}

	registry.InitByYaml(yamlConfig)

	// register entities
	registry.RegisterEntity(&entitybeeorm.NotificationEntity{})

	// register queues (redis streams)
	// dirty queues - every time entity is flushed, the ORM will push its ID in the queue

	registry.RegisterRedisStream(
		queue.OrmDirtyNotificationEntity,
		streamsPool,
		[]string{queue.GetConsumerGroupName(queue.OrmDirtyNotificationEntity)},
	)

	//orm internal queues
	//it's important for them to not be created in default redis pool, as every time cache is deleted, the consumer will panic

	//lazy flush - not used in current project, but still needs to be registered
	registry.RegisterRedisStream(beeorm.LazyChannelName, streamsPool, []string{beeorm.AsyncConsumerGroupName})
	// logs tables - not used in current project, but still needs to be registered
	registry.RegisterRedisStream(beeorm.LogChannelName, streamsPool, []string{beeorm.AsyncConsumerGroupName})
	// redis search indexer - not used in current project, but still needs to be registered
	registry.RegisterRedisStream(beeorm.RedisSearchIndexerChannelName, streamsPool, []string{beeorm.AsyncConsumerGroupName})
	//redis streams garbage collector - used
	registry.RegisterRedisStream(beeorm.RedisStreamGarbageCollectorChannelName, streamsPool, []string{beeorm.AsyncConsumerGroupName})

	validatedRegistry, defferFunc, err := registry.Validate()

	return &beeORMRegistry{validatedRegistry}, defferFunc, err
}

type beeORMRegistry struct {
	beeorm.ValidatedRegistry
}

func (b *beeORMRegistry) GetORMService() Engine {
	return &beeORMEngine{b.ValidatedRegistry.CreateEngine()}
}
