package main

import (
	"fmt"
	"syscall/js"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func main() {
	ch := make(chan struct{})
	js.Global().Set("format", js.FuncOf(format))
	<-ch
}

func format(js.Value, []js.Value) interface{} {
	input := js.Global().Get("document").Call("getElementById", "input").Get("value").String()
	output, err := j2s.Format(input)
	if err != nil {
		fmt.Println(err)
	}
	js.Global().Get("document").Call("getElementById", "output").Set("value", output)
	return nil
}
