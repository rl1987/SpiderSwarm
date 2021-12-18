package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
