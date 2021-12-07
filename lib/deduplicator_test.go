package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeduplicator(t *testing.T) {
	deduplicator := NewDeduplicator("localhost:6379")

	assert.NotNil(t, deduplicator)
	assert.NotNil(t, deduplicator.Backend)
}
