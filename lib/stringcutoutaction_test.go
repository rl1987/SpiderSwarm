package spiderswarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStringCutActionFromTemplate(t *testing.T) {
	constructorParams := map[string]interface{}{
		"from": "here",
		"to":   "eternity",
	}

	actionTempl := &ActionTemplate{
		StructName:        "StringCutAction",
		ConstructorParams: constructorParams,
	}

	action := NewStringCutActionFromTemplate(actionTempl)

	assert.NotNil(t, action)
	assert.Equal(t, "here", action.From)
	assert.Equal(t, "eternity", action.To)
}
