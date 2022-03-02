package main

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func main() {
	ch := make(chan struct{})
	js.Global().Set("json2struct", js.FuncOf(json2struct))
	js.Global().Set("copyClipboard", js.FuncOf(copyClipboard))
	<-ch
}

func json2struct(js.Value, []js.Value) interface{} {
	doc := js.Global().Get("document")
	input := doc.Call("getElementById", "input").Get("value").String()
	if input == "" {
		// if there is no input, output samples with placeholder value
		input = `{"sample":"paste the json here"}`
	}

	output, err := j2s.ConvertWithOption(input, j2s.Option{
		UseTag:    doc.Call("getElementById", "use-tag").Get("checked").Bool(),
		TagName:   doc.Call("getElementById", "tag-name").Get("value").String(),
		Omitempty: omitempty(&doc),
	})

	if err != nil {
		doc.Call("getElementById", "output").Set("value", err.Error())
		return nil
	}
	doc.Call("getElementById", "output").Set("value", output)
	return nil
}

func copyClipboard(js.Value, []js.Value) interface{} {
	output := js.Global().Get("document").Call("getElementById", "output").Get("value").String()
	js.Global().Get("navigator").Get("clipboard").Call("writeText", output)
	return nil
}

func omitempty(doc *js.Value) j2s.Omitempty {
	values := doc.Call("getElementsByName", "omitempty")
	for i := 0; i < values.Length(); i++ {
		idx := values.Index(i)
		if idx.Get("checked").Bool() {
			v, err := strconv.Atoi(idx.Get("value").String())
			if err != nil {
				fmt.Println(err)
				break
			}
			return j2s.Omitempty(v)
		}
	}
	return j2s.OmitemptyNone
}
