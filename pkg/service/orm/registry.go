package orm

import (
	"errors"
	"fmt"

	entitybeeorm "github.com/iliyanm/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanm/notification/pkg/queue"
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
	registry.RegisterEntity(&entitybeeorm.SMSNotificationEntity{})
	registry.RegisterEntity(&entitybeeorm.EmailNotificationEntity{})
	registry.RegisterEntity(&entitybeeorm.SlackNotificationEntity{})

	// register enums
	registry.RegisterEnumStruct("entitybeeorm.NotificationStatusAll", entitybeeorm.NotificationStatusAll)

	// register queues (redis streams)
	// dirty queues - every time entity is flushed, the ORM will push its ID in the queue

	registry.RegisterRedisStream(
		queue.OrmDirtyNotificationEntity,
		streamsPool,
		[]string{queue.GetConsumerGroupName(queue.OrmDirtyNotificationEntity)},
	)

	validatedRegistry, defferFunc, err := registry.Validate()

	return &beeORMRegistry{validatedRegistry}, defferFunc, err
}

type beeORMRegistry struct {
	beeorm.ValidatedRegistry
}

func (b *beeORMRegistry) GetORMService() Engine {
	return &beeORMEngine{b.ValidatedRegistry.CreateEngine()}
}
