package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRunner(t *testing.T) {
	runner := NewRunner("localhost:1337")

	assert.NotNil(t, runner)
	assert.Equal(t, "localhost:1337", runner.BackendAddr)
}
