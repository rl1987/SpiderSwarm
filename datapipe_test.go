package spiderswarm

import (
	"errors"
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

func TestDataPipeAddUnsupported(t *testing.T) {
	err := errors.New("Unsupported payload type")

	dataPipe := NewDataPipe()

	gotErr := dataPipe.Add(err)

	assert.NotNil(t, gotErr)
	assert.Equal(t, err, gotErr)
}
