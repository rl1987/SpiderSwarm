package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFormExtractionAction(t *testing.T) {
	action := NewFormExtractionAction("f1")

	assert.NotNil(t, action)
	assert.Equal(t, "f1", action.FormID)
	assert.Equal(t, []string{
		FormExtractionActionInputHTMLStr,
		FormExtractionActionInputHTMLBytes,
	}, action.AbstractAction.AllowedInputNames)
	assert.Equal(t, []string{
		FormExtractionActionOutputFormData,
	}, action.AbstractAction.AllowedOutputNames)
}

func TestNewFormExtractionActionFromTemplate(t *testing.T) {
	actionTempl := &ActionTemplate{
		Name:       "GetForm",
		StructName: "FormExtractionAction",
		ConstructorParams: map[string]Value{
			"formID": Value{
				ValueType:   ValueTypeString,
				StringValue: "f1",
			},
		},
	}

	action := NewFormExtractionActionFromTemplate(actionTempl, "").(*FormExtractionAction)

	assert.NotNil(t, action)
	assert.Equal(t, "GetForm", action.Name)
	assert.Equal(t, "f1", action.FormID)
}

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
		"custId":   "3487",
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
