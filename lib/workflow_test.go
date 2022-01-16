package spsw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddActionTemplate(t *testing.T) {
	taskTempl := NewTaskTemplate("testTask", false)

	expectTaskTempl := &TaskTemplate{
		TaskName: "testTask",
		Initial:  false,
		ActionTemplates: []ActionTemplate{
			ActionTemplate{
				Name:              "HTTP",
				StructName:        "HTTPAction",
				ConstructorParams: map[string]Value{},
			},
			ActionTemplate{
				Name:              "XPath",
				StructName:        "XPathAction",
				ConstructorParams: map[string]Value{},
			},
			ActionTemplate{
				Name:              "Join",
				StructName:        "FieldJoinAction",
				ConstructorParams: map[string]Value{},
			},
		},
		DataPipeTemplates: []DataPipeTemplate{},
	}

	taskTempl.AddActionTemplate(NewActionTemplate("HTTP", "HTTPAction", nil))
	taskTempl.AddActionTemplate(NewActionTemplate("XPath", "XPathAction", map[string]interface{}{}))
	taskTempl.AddActionTemplate(NewActionTemplate("Join", "FieldJoinAction", map[string]interface{}{}))

	err := taskTempl.AddActionTemplate(NewActionTemplate("HTTP", "HTTPAction", nil))
	assert.NotNil(t, err)

	assert.Equal(t, expectTaskTempl, taskTempl)
}

func TestRemoveActionTemplate(t *testing.T) {
	taskTempl := &TaskTemplate{
		TaskName: "testTask",
		Initial:  false,
		ActionTemplates: []ActionTemplate{
			ActionTemplate{
				Name:              "HTTP",
				StructName:        "HTTPAction",
				ConstructorParams: map[string]Value{},
			},
			ActionTemplate{
				Name:              "XPath",
				StructName:        "XPathAction",
				ConstructorParams: map[string]Value{},
			},
			ActionTemplate{
				Name:              "Join",
				StructName:        "FieldJoinAction",
				ConstructorParams: map[string]Value{},
			},
		},
		DataPipeTemplates: []DataPipeTemplate{},
	}

	err := taskTempl.RemoveActionTemplate("!")
	assert.NotNil(t, err)

	err = taskTempl.RemoveActionTemplate("HTTP")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(taskTempl.ActionTemplates))

	err = taskTempl.RemoveActionTemplate("XPath")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(taskTempl.ActionTemplates))

	err = taskTempl.RemoveActionTemplate("Join")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(taskTempl.ActionTemplates))
}

func TestWorkflowYAMLAndBack(t *testing.T) {
	workflow := &Workflow{
		Name:    "testWorkflow",
		Version: "v0.0.0.0.1",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPAction",
						ConstructorParams: map[string]Value{
							"baseURL": Value{
								ValueType:   ValueTypeString,
								StringValue: "https://news.ycombinator.com/",
							},
							"method": Value{
								ValueType:   ValueTypeString,
								StringValue: "GET",
							},
							"canFail": Value{
								ValueType: ValueTypeBool,
								BoolValue: false,
							},
						},
					},
					ActionTemplate{
						Name:              "UTF8Decode",
						StructName:        "UTF8DecodeAction",
						ConstructorParams: map[string]Value{},
					},
					ActionTemplate{
						Name:       "MakePromise",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]Value{
							"inputNames": Value{
								ValueType:    ValueTypeStrings,
								StringsValue: []string{"htmlStr1", "htmlStr2"},
							},
							"taskName": Value{
								ValueType:   ValueTypeString,
								StringValue: "ParseHTML",
							},
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						SourceActionName: "HTTP1",
						SourceOutputName: HTTPActionOutputBody,
						DestActionName:   "UTF8Decode",
						DestInputName:    UTF8DecodeActionInputBytes,
					},
					DataPipeTemplate{
						SourceActionName: "UTF8Decode",
						SourceOutputName: UTF8DecodeActionOutputStr,
						DestActionName:   "MakePromise",
						DestInputName:    "htmlStr1",
					},
					DataPipeTemplate{
						SourceActionName: "UTF8Decode",
						SourceOutputName: UTF8DecodeActionOutputStr,
						DestActionName:   "MakePromise",
						DestInputName:    "htmlStr2",
					},
					DataPipeTemplate{
						SourceActionName: "MakePromise",
						SourceOutputName: TaskPromiseActionOutputPromise,
						TaskOutputName:   "promise",
					},
				},
			},
			TaskTemplate{
				TaskName: "ParseHTML",
				Initial:  false,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "TitleExtraction",
						StructName: "XPathAction",
						ConstructorParams: map[string]Value{
							"xpath": Value{
								ValueType:   ValueTypeString,
								StringValue: "//a[@class='storylink']/text()",
							},
							"expectMany": Value{
								ValueType: ValueTypeBool,
								BoolValue: true,
							},
						},
					},
					ActionTemplate{
						Name:       "LinkExtraction",
						StructName: "XPathAction",
						ConstructorParams: map[string]Value{
							"xpath": Value{
								ValueType:   ValueTypeString,
								StringValue: "//a[@class='storylink']/@href",
							},
							"expectMany": Value{
								ValueType: ValueTypeBool,
								BoolValue: true,
							},
						},
					},
					ActionTemplate{
						Name:       "YieldItem",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]Value{
							"inputNames": Value{
								ValueType:    ValueTypeStrings,
								StringsValue: []string{"title", "link"},
							},
							"itemName": Value{
								ValueType:   ValueTypeString,
								StringValue: "HNItem",
							},
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "htmlStr1",
						DestActionName: "TitleExtraction",
						DestInputName:  XPathActionInputHTMLStr,
					},
					DataPipeTemplate{
						TaskInputName:  "htmlStr2",
						DestActionName: "LinkExtraction",
						DestInputName:  XPathActionInputHTMLStr,
					},
					DataPipeTemplate{
						SourceActionName: "TitleExtraction",
						SourceOutputName: XPathActionOutputStr,
						DestActionName:   "YieldItem",
						DestInputName:    "title",
					},
					DataPipeTemplate{
						SourceActionName: "LinkExtraction",
						SourceOutputName: XPathActionOutputStr,
						DestActionName:   "YieldItem",
						DestInputName:    "link",
					},
					DataPipeTemplate{
						SourceActionName: "YieldItem",
						SourceOutputName: FieldJoinActionOutputItem,
						TaskOutputName:   "items",
					},
				},
			},
		},
	}

	yamlStr := workflow.ToYAML()
	gotWorkflow := NewWorkflowFromYAML(yamlStr)

	assert.Equal(t, workflow, gotWorkflow)
}

func TestWorkflowValidateActionStructNames(t *testing.T) {
	workflow := &Workflow{
		Name:    "testWorkflow",
		Version: "v0.0.0.0.1",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPActionn", // Typo!
					},
				},
			},
		},
	}

	err := workflow.validateActionStructNames()
	assert.NotNil(t, err)

	workflow.TaskTemplates[0].ActionTemplates[0].StructName = "HTTPAction"

	err = workflow.validateActionStructNames()
	assert.Nil(t, err)
}

func TestWorkflowValidateActionConnectedness(t *testing.T) {
	workflow := &Workflow{
		Name:    "testWorkflow",
		Version: "v0.0.0.0.1",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPAction",
					},
					ActionTemplate{
						Name:       "Unconnected",
						StructName: "ConstAction",
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "url_params",
						DestActionName: "HTTP1",
						DestInputName:  HTTPActionInputURLParams,
					},
					DataPipeTemplate{
						SourceActionName: "HTTP1",
						SourceOutputName: HTTPActionOutputBody,
						TaskOutputName:   "body",
					},
				},
			},
		},
	}

	err := workflow.validateActionConnectedness()

	assert.NotNil(t, err)

	workflow.TaskTemplates[0].ActionTemplates = []ActionTemplate{
		workflow.TaskTemplates[0].ActionTemplates[0],
	}

	err = workflow.validateActionConnectedness()

	assert.Nil(t, err)
}

func TestWorkflowValidateInputOutputNames(t *testing.T) {
	workflow1 := &Workflow{
		Name:    "testWorkflow1",
		Version: "v0.0.0.0.1",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPAction",
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "url_params",
						DestActionName: "HTTP1",
						DestInputName:  "params",
					},
				},
			},
		},
	}

	err := workflow1.validateInputOutputNames()
	assert.NotNil(t, err)

	workflow2 := &Workflow{
		Name:    "testWorkflow2",
		Version: "v0.0.0.0.1",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPAction",
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "url_params",
						DestActionName: "HTTP1",
						DestInputName:  HTTPActionInputURLParams,
					},
					DataPipeTemplate{
						SourceActionName: "HTTTP1",
						SourceOutputName: "body",
						TaskOutputName:   "body",
					},
				},
			},
		},
	}

	err = workflow2.validateInputOutputNames()
	assert.NotNil(t, err)

	workflow3 := &Workflow{
		Name:    "testWorkflow3",
		Version: "v0.0.0.0.1",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPAction",
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "url_params",
						DestActionName: "HTTP1",
						DestInputName:  HTTPActionInputURLParams,
					},
					DataPipeTemplate{
						SourceActionName: "HTTP1",
						SourceOutputName: HTTPActionOutputBody,
						TaskOutputName:   "body",
					},
				},
			},
		},
	}

	err = workflow3.validateInputOutputNames()
	assert.Nil(t, err)
}

func TestNewWorkflow(t *testing.T) {
	expectWorkflow := &Workflow{
		Name:          "Test",
		Version:       "1",
		TaskTemplates: []TaskTemplate{},
	}

	gotWorkflow := NewWorkflow("Test", "1")

	assert.Equal(t, expectWorkflow, gotWorkflow)
}

func TestWorkflowAddTaskTemplate(t *testing.T) {
	workflow := &Workflow{}
	taskTempl := &TaskTemplate{TaskName: "TestTask"}

	workflow.AddTaskTemplate(taskTempl)

	assert.Equal(t, 1, len(workflow.TaskTemplates))
	assert.Equal(t, taskTempl, &workflow.TaskTemplates[0])
}

func TestWorkflowSetInitial(t *testing.T) {
	workflow := &Workflow{
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "First",
				Initial:  false,
			},
			TaskTemplate{
				TaskName: "Second",
				Initial:  false,
			},
		},
	}

	workflow.SetInitial("First")

	assert.True(t, workflow.TaskTemplates[0].Initial)
	assert.False(t, workflow.TaskTemplates[1].Initial)

	workflow.SetInitial("Second")

	assert.False(t, workflow.TaskTemplates[0].Initial)
	assert.True(t, workflow.TaskTemplates[1].Initial)
}

func TestWorkflowGetInitialTaskTemplate(t *testing.T) {
	workflow := &Workflow{
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "First",
				Initial:  false,
			},
			TaskTemplate{
				TaskName: "Second",
				Initial:  true,
			},
		},
	}

	initialTT := workflow.GetInitialTaskTemplate()

	assert.Equal(t, &workflow.TaskTemplates[1], initialTT)
}

func TestWorkflowRemoveTaskTemplate(t *testing.T) {
	workflow := &Workflow{
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "First",
				Initial:  true,
			},
			TaskTemplate{
				TaskName: "Second",
				Initial:  false,
			},
			TaskTemplate{
				TaskName: "Third",
				Initial:  false,
			},
		},
	}

	err := workflow.RemoveTaskTemplate("Second")

	assert.Nil(t, err)

	assert.Equal(t, 2, len(workflow.TaskTemplates))
	assert.Equal(t, "First", workflow.TaskTemplates[0].TaskName)
	assert.Equal(t, "Third", workflow.TaskTemplates[1].TaskName)

	err = workflow.RemoveTaskTemplate("Fourth")
	assert.NotNil(t, err)
}

func TestWorkflowGetInitialTaskTemplateName(t *testing.T) {
	workflow := &Workflow{
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "First",
				Initial:  true,
			},
			TaskTemplate{
				TaskName: "Second",
				Initial:  false,
			},
			TaskTemplate{
				TaskName: "Third",
				Initial:  false,
			},
		},
	}

	taskTemplName := workflow.GetInitialTaskTemplateName()

	assert.Equal(t, "First", taskTemplName)
}

func TestTaskTemplateConnectActionTemplates(t *testing.T) {
	tt := &TaskTemplate{
		DataPipeTemplates: []DataPipeTemplate{},
	}

	expectDataPipeTemplate := DataPipeTemplate{
		SourceActionName: "Action1",
		SourceOutputName: "out1",
		DestActionName:   "Action2",
		DestInputName:    "in2",
	}

	err := tt.ConnectActionTemplates("Action1", "out1", "Action2", "in2")

	assert.Nil(t, err)

	assert.Equal(t, 1, len(tt.DataPipeTemplates))
	assert.Equal(t, expectDataPipeTemplate, tt.DataPipeTemplates[0])
}

func TestTaskTemplateConnectInputToActionTemplate(t *testing.T) {
	tt := &TaskTemplate{
		DataPipeTemplates: []DataPipeTemplate{},
	}

	expectDataPipeTemplate := DataPipeTemplate{
		TaskInputName:  "in1",
		DestActionName: "Action1",
		DestInputName:  "in2",
	}

	err := tt.ConnectInputToActionTemplate("in1", "Action1", "in2")

	assert.Nil(t, err)

	assert.Equal(t, 1, len(tt.DataPipeTemplates))
	assert.Equal(t, expectDataPipeTemplate, tt.DataPipeTemplates[0])
}

func TestTaskTemplateConnectOutputToActionTemplate(t *testing.T) {
	tt := &TaskTemplate{
		DataPipeTemplates: []DataPipeTemplate{},
	}

	expectDataPipeTemplate := DataPipeTemplate{
		SourceActionName: "Action1",
		SourceOutputName: "out1",
		TaskOutputName:   "out2",
	}

	err := tt.ConnectOutputToActionTemplate("Action1", "out1", "out2")

	assert.Nil(t, err)

	assert.Equal(t, 1, len(tt.DataPipeTemplates))
	assert.Equal(t, expectDataPipeTemplate, tt.DataPipeTemplates[0])
}

func TestTaskTemplateDisconnectActionTemplates(t *testing.T) {
	tt := TaskTemplate{
		DataPipeTemplates: []DataPipeTemplate{
			DataPipeTemplate{
				SourceActionName: "Action1",
				SourceOutputName: "out1",
				DestActionName:   "Action2",
				DestInputName:    "in2",
			},
		},
	}

	err := tt.DisconnectActionTemplates("Action1", "out1", "Action2", "in2")

	assert.Nil(t, err)

	assert.Equal(t, 0, len(tt.DataPipeTemplates))

	err = tt.DisconnectActionTemplates("Action2", "out3", "Action3", "in4")
	assert.Equal(t, errors.New("Not found"), err)
}

func TestDisconnectInput(t *testing.T) {
	tt := TaskTemplate{
		DataPipeTemplates: []DataPipeTemplate{
			DataPipeTemplate{
				TaskInputName:  "in1",
				DestActionName: "Action1",
				DestInputName:  "in2",
			},
		},
	}

	err := tt.DisconnectInput("in1")

	assert.Nil(t, err)

	assert.Equal(t, 0, len(tt.DataPipeTemplates))

	err = tt.DisconnectInput("in2")
	assert.Equal(t, errors.New("Not found"), err)
}

func TestDisconnectOutput(t *testing.T) {
	tt := TaskTemplate{
		DataPipeTemplates: []DataPipeTemplate{
			DataPipeTemplate{
				SourceActionName: "Action1",
				SourceOutputName: "out1",
				TaskOutputName:   "out2",
			},
		},
	}

	err := tt.DisconnectOutput("out2")

	assert.Nil(t, err)

	assert.Equal(t, 0, len(tt.DataPipeTemplates))

	err = tt.DisconnectOutput("out4")
	assert.Equal(t, errors.New("Not found"), err)
}
