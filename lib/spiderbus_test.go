package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpiderBus(t *testing.T) {
	spiderBus := NewSpiderBus()
	assert.NotNil(t, spiderBus)
	assert.NotNil(t, spiderBus.UUID)
}
