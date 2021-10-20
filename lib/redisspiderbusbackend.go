package spsw

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisSpiderBusBackend struct {
	SpiderBusBackend

	ctx         context.Context
	serverAddr  string
	redisClient *redis.Client
	consumerId  string
}

func NewRedisSpiderBusBackend(serverAddr string, password string) *RedisSpiderBusBackend {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     serverAddr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()

	consumerId := uuid.New().String()

	status := redisClient.XGroupCreateMkStream(ctx, "items", consumerId, "$")
	spew.Dump(status)
	status = redisClient.XGroupCreateMkStream(ctx, "task_promises", consumerId, "$")
	spew.Dump(status)
	status = redisClient.XGroupCreateMkStream(ctx, "scheduled_tasks", consumerId, "$")
	spew.Dump(status)

	return &RedisSpiderBusBackend{
		ctx:         ctx,
		serverAddr:  serverAddr,
		redisClient: redisClient,
		consumerId:  consumerId,
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
	rsbb.redisClient.Close()
}
