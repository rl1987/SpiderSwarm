package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataPipeBetweenAction(t *testing.T) {
	httpAction := &HTTPAction{}
	xpathAction := &XPathAction{}

	dataPipe := NewDataPipeBetweenActions(httpAction, xpathAction)

	assert.NotNil(t, dataPipe)
	assert.Equal(t, httpAction, dataPipe.FromAction)
	assert.Equal(t, xpathAction, dataPipe.ToAction)
	assert.NotNil(t, dataPipe.Queue)
}
