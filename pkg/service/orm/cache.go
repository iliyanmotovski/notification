package orm

import "github.com/latolukasz/beeorm"

type CacheService interface {
	Set(key string, values ...interface{})
	Get(key, field string) (string, bool)
}

type redisCache struct {
	*beeorm.RedisCache
}

func (r *redisCache) Set(key string, values ...interface{}) {
	r.RedisCache.HSet(key, values...)
}

func (r *redisCache) Get(key, field string) (string, bool) {
	return r.RedisCache.HGet(key, field)
}
