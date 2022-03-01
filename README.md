# go-json2struct
a web tool for converting JSON strings into Go struct

The web application is available at [here](https://masakurapa.github.io/go-json2struct/) !!

## Example for use as a module

Check [pkg.go.dev](https://pkg.go.dev/github.com/masakurapa/go-json2struct/pkg/j2s) for examples.

```sh
$ go get -u github.com/masakurapa/go-json2struct
```

```go
package main

import (
	"fmt"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func main() {
	input := `{"title": "j2s"}`
	output, err := j2s.Convert(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output)
}
```
