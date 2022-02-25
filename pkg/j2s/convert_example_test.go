package j2s_test

import (
	"fmt"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func ExampleConvert() {
	// TODO: The order cannot be fixed because it uses a map
	// input := `{
	// 	"title": "j2s",
	// 	"snake_case": 99,
	// 	"CamelCase": true,
	// 	"kebab-case": null,
	// 	"map": {"child1": "apple", "child2": 12345},
	// 	"slice": ["1", "2", "3", "4", "5"]
	// }`

	input := `{"title": "j2s"}`

	output, err := j2s.Convert(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output)

	// Output: type J2S1 struct {
	// 	Title string `json:"title"`
	// }
}
