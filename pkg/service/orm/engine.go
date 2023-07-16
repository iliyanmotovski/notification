package orm

import (
	"log"

	"github.com/latolukasz/beeorm"
)

type Engine interface {
	Flush(entity interface{})
	NewFlusher() Flusher
	LoadByID(id uint64, entity interface{}, references ...string) bool
	ExecuteAlters()
	TruncateTables()
	GetEventBroker() EventBroker
	GetCacheService(namespace ...string) CacheService
	GetRegistry() RegistryService
}

type beeORMEngine struct {
	*beeorm.Engine
}

func (b *beeORMEngine) Flush(entity interface{}) {
	b.Engine.Flush(entity.(beeorm.Entity))
}

func (b *beeORMEngine) NewFlusher() Flusher {
	return &beeORMFlusher{b.Engine.NewFlusher()}
}

func (b *beeORMEngine) LoadByID(id uint64, entity interface{}, references ...string) bool {
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

func (b *beeORMEngine) TruncateTables() {
	dbService := b.Engine.GetMysql()

	var query string
	rows, deferF := dbService.Query(
		"SELECT CONCAT('delete from  ',table_schema,'.',table_name,';' , 'ALTER TABLE ', table_schema,'.',table_name , ' AUTO_INCREMENT = 1;') AS query " +
			"FROM information_schema.tables WHERE table_schema IN ('" + dbService.GetPoolConfig().GetDatabase() + "');",
	)

	defer deferF()

	if rows != nil {
		var queries string

		for rows.Next() {
			rows.Scan(&query)
			queries += query
		}

		_, def := dbService.Query("SET FOREIGN_KEY_CHECKS=0;" + queries + "SET FOREIGN_KEY_CHECKS=1")
		defer def()
	}
}

func (b *beeORMEngine) GetEventBroker() EventBroker {
	return &beeORMEventBroker{b.Engine.GetEventBroker()}
}

func (b *beeORMEngine) GetCacheService(namespace ...string) CacheService {
	return &redisCache{b.Engine.GetRedis(namespace...)}
}

func (b *beeORMEngine) GetRegistry() RegistryService {
	return &beeORMRegistry{b.Engine.GetRegistry()}
}
