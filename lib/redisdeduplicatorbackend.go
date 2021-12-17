package spsw

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type RedisDeduplicatorBackend struct {
	AbstractDeduplicatorBackend
	UUID string

	ctx         context.Context
	serverAddr  string
	redisClient *redis.Client
}

func NewRedisDeduplicatorBackend(serverAddr string, password string) *RedisDeduplicatorBackend {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     serverAddr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()

	return &RedisDeduplicatorBackend{
		UUID:        uuid.New().String(),
		ctx:         ctx,
		serverAddr:  serverAddr,
		redisClient: redisClient,
	}
}

func (rdb *RedisDeduplicatorBackend) isHashInRedisSet(key string, hashStr string) bool {
	resp := rdb.redisClient.SIsMember(rdb.ctx, key, hashStr)
	res, err := resp.Result()
	if err != nil {
		log.Error(fmt.Sprintf("Redis SISMEMBER command failed with error: %v", err))
		return false
	}

	return res
}

func (rdb *RedisDeduplicatorBackend) IsScheduledTaskDuplicated(scheduledTask *ScheduledTask) bool {
	hashBytes := scheduledTask.Hash()
	hashStr := fmt.Sprintf("%v", hashBytes)
	key := "scheduledtasks-" + scheduledTask.JobUUID

	return rdb.isHashInRedisSet(key, hashStr)
}

func (rdb *RedisDeduplicatorBackend) NoteScheduledTask(scheduledTask *ScheduledTask) error {
	hashBytes := scheduledTask.Hash()
	hashStr := fmt.Sprintf("%v", hashBytes)
	key := "scheduledtasks-" + scheduledTask.JobUUID

	resp := rdb.redisClient.SAdd(rdb.ctx, key, hashStr)
	return resp.Err()
}

func (rdb *RedisDeduplicatorBackend) Close() {
	rdb.redisClient.Close()
}
