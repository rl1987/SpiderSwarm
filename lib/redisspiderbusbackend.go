package spsw

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisSpiderBusBackend struct {
	SpiderBusBackend

	ctx         context.Context
	serverAddr  string
	redisClient *redis.Client
}

func NewRedisSpiderBusBackend(serverAddr string, password string) *RedisSpiderBusBackend {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     serverAddr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()

	return &RedisSpiderBusBackend{
		ctx:         ctx,
		serverAddr:  serverAddr,
		redisClient: redisClient,
	}
}

func (rsbb *RedisSpiderBusBackend) SendScheduledTask(scheduledTask *ScheduledTask) error {
	return nil
}

func (rsbb *RedisSpiderBusBackend) ReceiveScheduledTask() *ScheduledTask {
	return nil
}

func (rsbb *RedisSpiderBusBackend) SendTaskPromise(taskPromise *TaskPromise) error {
	return nil
}

func (rsbb *RedisSpiderBusBackend) ReceiveTaskPromise() *TaskPromise {
	return nil
}

func (rsbb *RedisSpiderBusBackend) SendItem(item *Item) error {
	return nil
}

func (rsbb *RedisSpiderBusBackend) ReceiveItem() *Item {
	return nil
}

func (rsbb *RedisSpiderBusBackend) Close() {
}
