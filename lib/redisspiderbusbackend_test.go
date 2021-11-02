package spsw

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedisSpiderBusBackend(t *testing.T) {
	rsbb := NewRedisSpiderBusBackend("127.0.0.1:6379", "")

	assert.NotNil(t, rsbb)
	assert.NotNil(t, rsbb.redisClient)
	assert.Equal(t, "127.0.0.1:6379", rsbb.serverAddr)
	assert.Equal(t, context.Background(), rsbb.ctx)
}
