package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSplayable(t *testing.T) {
	item1 := &Item{
		Fields: map[string]interface{}{
			"field1": "a",
			"field2": "b",
		},
	}

	assert.False(t, item1.IsSplayable())

	item2 := &Item{
		Fields: map[string]interface{}{
			"field1": []string{"a", "b", "c"},
			"field2": []string{"1", "2", "3"},
		},
	}

	assert.True(t, item2.IsSplayable())

	item3 := &Item{
		Fields: map[string]interface{}{
			"field1": []string{"x", "y", "z"},
			"field2": []string{"0", "1"},
		},
	}

	assert.True(t, item3.IsSplayable())
}
