package spsw

import (
	"context"
	"errors"
	"time"

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

	redisClient.XGroupCreateMkStream(ctx, "items", consumerId, "$")
	redisClient.XGroupCreateMkStream(ctx, "task_promises", consumerId, "$")
	redisClient.XGroupCreateMkStream(ctx, "scheduled_tasks", consumerId, "$")

	return &RedisSpiderBusBackend{
		ctx:         ctx,
		serverAddr:  serverAddr,
		redisClient: redisClient,
		consumerId:  consumerId,
	}
}

func (rsbb *RedisSpiderBusBackend) SendScheduledTask(scheduledTask *ScheduledTask) error {
	raw := scheduledTask.EncodeToJSON()

	resp := rsbb.redisClient.XAdd(rsbb.ctx, &redis.XAddArgs{
		Stream: "scheduled_tasks",
		ID:     "*",
		Values: map[string]interface{}{
			"raw": string(raw),
		},
	})

	spew.Dump(resp)

	err := resp.Err()

	if err != nil {
		spew.Dump(err)
		return err
	}

	return nil
}

func (rsbb *RedisSpiderBusBackend) readRawMessageFromStream(stream string) ([]byte, error) {
	resp := rsbb.redisClient.XReadGroup(rsbb.ctx, &redis.XReadGroupArgs{
		Group:    rsbb.consumerId,
		Consumer: rsbb.consumerId,
		Streams:  []string{stream, ">"},
		Count:    1,
		Block:    1 * time.Second,
		NoAck:    false,
	})

	s, err := resp.Result()
	if err != nil {
		return nil, err
	}

	if len(s) == 1 && len(s[0].Messages) == 1 {
		msg := s[0].Messages[0]
		if raw, ok := msg.Values["raw"].(string); ok {
			rsbb.redisClient.XAck(rsbb.ctx, stream, rsbb.consumerId, msg.ID)
			return []byte(raw), nil
		}
	}

	return nil, errors.New("Unknown error")
}

func (rsbb *RedisSpiderBusBackend) ReceiveScheduledTask() *ScheduledTask {
	raw, err := rsbb.readRawMessageFromStream("scheduled_tasks")
	if err != nil {
		return nil
	}

	scheduledTask := NewScheduledTaskFromJSON(raw)

	return scheduledTask
}

func (rsbb *RedisSpiderBusBackend) SendTaskPromise(taskPromise *TaskPromise) error {
	raw := taskPromise.EncodeToJSON()

	err := rsbb.redisClient.XAdd(rsbb.ctx, &redis.XAddArgs{
		Stream: "task_promises",
		ID:     "$",
		Values: map[string]interface{}{
			"raw": string(raw),
		},
	}).Err()

	if err != nil {
		spew.Dump(err)
		return err
	}

	return nil
}

func (rsbb *RedisSpiderBusBackend) ReceiveTaskPromise() *TaskPromise {
	raw, err := rsbb.readRawMessageFromStream("task_promises")
	if err != nil {
		return nil
	}

	taskPromise := NewTaskPromiseFromJSON(raw)

	return taskPromise
}

func (rsbb *RedisSpiderBusBackend) SendItem(item *Item) error {
	raw := item.EncodeToJSON()

	err := rsbb.redisClient.XAdd(rsbb.ctx, &redis.XAddArgs{
		Stream: "items",
		ID:     "$",
		Values: map[string]interface{}{
			"raw": string(raw),
		},
	}).Err()

	if err != nil {
		spew.Dump(err)
		return err
	}

	return nil
}

func (rsbb *RedisSpiderBusBackend) ReceiveItem() *Item {
	raw, err := rsbb.readRawMessageFromStream("items")
	if err != nil {
		return nil
	}

	item := NewItemFromJSON(raw)

	return item
}

func (rsbb *RedisSpiderBusBackend) Close() {
	for _, stream := range []string{"items", "scheduled_tasks", "task_promises"} {
		rsbb.redisClient.XGroupDelConsumer(rsbb.ctx, stream, rsbb.consumerId, rsbb.consumerId)
		rsbb.redisClient.XGroupDestroy(rsbb.ctx, stream, rsbb.consumerId)
	}

	rsbb.redisClient.Close()
}
