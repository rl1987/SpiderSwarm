package main

import (
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	spsw "github.com/spiderswarm/spiderswarm/lib"
)

func main() {
	taskTempl := spsw.NewTaskTemplate("TestTask", false)

	httpActionTempl := spsw.NewActionTemplate("GetHTML", "HTTPAction", map[string]interface{}{
		"baseURL": "https://quotes.toscrape.com/",
		"method":  "GET",
		"canFail": false,
	})
	taskTempl.AddActionTemplate(httpActionTempl)

	xpathActionTempl1 := spsw.NewActionTemplate("XPathQuote", "XPathAction", map[string]interface{}{
		"xpath":      "//div[@class=\"quote\"]/span",
		"expectMany": true,
	})
	taskTempl.AddActionTemplate(xpathActionTempl1)

	xpathActionTempl2 := spsw.NewActionTemplate("XPathAuthor", "XPathAction", map[string]interface{}{
		"xpath":      "//small[@class=\"author\"]/text()",
		"expectMany": true,
	})
	taskTempl.AddActionTemplate(xpathActionTempl2)

	fjaTempl := spsw.NewActionTemplate("MakeItem", "FieldJoinAction", map[string]interface{}{
		"inputNames": []string{"quote", "author"},
		"itemName":   "quote",
	})
	taskTempl.AddActionTemplate(fjaTempl)

	taskTempl.ConnectActionTemplates("GetHTML", spsw.HTTPActionOutputBody,
		"XPathQuote", spsw.XPathActionInputHTMLBytes)
	taskTempl.ConnectActionTemplates("GetHTML", spsw.HTTPActionOutputBody,
		"XPathAuthor", spsw.XPathActionInputHTMLBytes)
	taskTempl.ConnectActionTemplates("XPathQuote", spsw.XPathActionOutputStr,
		"MakeItem", "quote")
	taskTempl.ConnectActionTemplates("XPathAuthor", spsw.XPathActionOutputStr,
		"MakeItem", "author")

	taskTempl.ConnectOutputToActionTemplate("MakeItem", spsw.FieldJoinActionOutputItem, "quoteItem")

	spew.Dump(taskTempl)

	workflow := spsw.NewWorkflow("TestWF", "1")

	workflow.AddTaskTemplate(taskTempl)
	workflow.SetInitial("TestTask")

	_, err := workflow.Validate()
	if err != nil {
		spew.Dump(err)
	}

	spew.Dump(workflow)

	yamlStr := workflow.ToYAML()

	ioutil.WriteFile("test.yaml", []byte(yamlStr), 0644)
}
