package main

import (
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
	input := js.Global().Get("document").Call("getElementById", "input").Get("value").String()
	if input == "" {
		// if there is no input, output samples with placeholder value
		input = `{"sample":"paste the json here"}`
	}

	output, err := j2s.ConvertWithOption(input, j2s.Option{
		UseTag: true,
	})

	if err != nil {
		js.Global().Get("document").Call("getElementById", "output").Set("value", err.Error())
		return nil
	}
	js.Global().Get("document").Call("getElementById", "output").Set("value", output)
	return nil
}

func copyClipboard(js.Value, []js.Value) interface{} {
	output := js.Global().Get("document").Call("getElementById", "output").Get("value").String()
	js.Global().Get("navigator").Get("clipboard").Call("writeText", output)
	return nil
}
