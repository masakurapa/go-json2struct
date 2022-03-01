package j2s_test

import (
	"fmt"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func ExampleConvert() {
	input := `{
		"title": "j2s",
		"snake_case": 99,
		"CamelCase": true,
		"kebab-case": null,
		"map": {"child1": "apple", "child2": 12345},
		"array": ["1", "2", "3", "4", "5"]
	}`

	output, err := j2s.Convert(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output)

	// Output: type J2S1 struct {
	// 	CamelCase bool        `json:"CamelCase"`
	// 	Array     []string    `json:"array"`
	// 	KebabCase interface{} `json:"kebab-case"`
	// 	Map       J2S2        `json:"map"`
	// 	SnakeCase int         `json:"snake_case"`
	// 	Title     string      `json:"title"`
	// }
	//
	// type J2S2 struct {
	// 	Child1 string `json:"child1"`
	// 	Child2 int    `json:"child2"`
	// }
}
