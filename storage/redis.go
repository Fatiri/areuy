package storage

import (
	"github.com/go-redis/redis/v8"
)

type Redis interface {
	Run() *redis.Client
}

type RedisCtx struct {
	host, password string
	DB             int
}

func NewRedis(host, password string, DB int) Redis {
	return &RedisCtx{
		host:     host,
		password: password,
		DB:       DB,
	}
}

func (r *RedisCtx) Run() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     r.host,
		Password: r.password,
		DB:       r.DB,
	})
	return client
}
