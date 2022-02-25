package main

import (
	"syscall/js"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func main() {
	ch := make(chan struct{})
	js.Global().Set("json2struct", js.FuncOf(json2struct))
	<-ch
}

func json2struct(js.Value, []js.Value) interface{} {
	input := js.Global().Get("document").Call("getElementById", "input").Get("value").String()
	output, err := j2s.Convert(input)
	if err != nil {
		js.Global().Get("document").Call("getElementById", "output").Set("value", err.Error())
		return nil
	}
	js.Global().Get("document").Call("getElementById", "output").Set("value", output)
	return nil
}
