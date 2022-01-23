package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormExtractionActionRun(t *testing.T) {
	htmlStr := `
<html>
  <body>
    <form id="test">
       <input type="hidden" id="custId" name="custId" value="3487">
       <input type="hidden" id="custName" name="custName" value="John">
    </form>
  </body>
</html>
`
	expectFormData := map[string]string{
		"custId": "3487",
		"custName": "John",
	}

	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add(htmlStr)

	action := NewFormExtractionAction("test")

	action.AddInput(FormExtractionActionInputHTMLStr, dataPipeIn)
	action.AddOutput(FormExtractionActionOutputFormData, dataPipeOut)

	err := action.Run()
	assert.Nil(t, err)

	formData, ok := dataPipeOut.Remove().(map[string]string)
	assert.True(t, ok)

	assert.Equal(t, expectFormData, formData)
}

