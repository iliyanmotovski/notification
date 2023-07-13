package orm

import (
	"log"

	"github.com/latolukasz/beeorm"
)

type Engine interface {
	Flush(entity interface{})
	LoadFromCacheByID(id uint64, entity interface{}, references ...string) bool
	ExecuteAlters()
	GetEventBroker() EventBroker
	GetCacheService(namespace string) CacheService
}

type beeORMEngine struct {
	*beeorm.Engine
}

func (b *beeORMEngine) Flush(entity interface{}) {
	b.Engine.Flush(entity.(beeorm.Entity))
}

func (b *beeORMEngine) LoadFromCacheByID(id uint64, entity interface{}, references ...string) bool {
	return b.Engine.LoadByID(id, entity.(beeorm.Entity), references...)
}

func (b *beeORMEngine) ExecuteAlters() {
	hasAlters := false
	alters := b.Engine.GetAlters()

	for _, alter := range alters {
		hasAlters = true
		log.Println(alter.SQL)
		alter.Exec()
	}

	if hasAlters {
		b.Engine.GetRedis().FlushDB()
	}
}

func (b *beeORMEngine) GetEventBroker() EventBroker {
	return &beeORMEventBroker{b.Engine.GetEventBroker()}
}

func (b *beeORMEngine) GetCacheService(namespace string) CacheService {
	return &redisCache{b.Engine.GetRedis(namespace)}
}
