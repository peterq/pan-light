package cache

import (
	"encoding/json"
	"github.com/go-redis/cache"
	"github.com/peterq/pan-light/server/conf"
	"time"
)

var codec = &cache.Codec{
	Redis: conf.Redis,

	Marshal: func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	},
	Unmarshal: func(b []byte, v interface{}) error {
		return json.Unmarshal(b, v)
	},
}

func RedisGet(key string, data interface{}) error {
	return codec.Get(key, data)
}

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	return codec.Set(&cache.Item{
		Key:        key,
		Object:     value,
		Expiration: expiration,
	})
}
