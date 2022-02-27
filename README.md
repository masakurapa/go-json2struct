# go-json2struct
a web tool for converting JSON strings into Go struct

The web application is available at [here](https://masakurapa.github.io/go-json2struct/) !!

## Example

```go
func ExampleConvert() {
	input := `{
		"title": "j2s",
		"snake_case": 99,
		"CamelCase": true,
		"kebab-case": null,
		"map": {"child1": "apple", "child2": 12345},
		"slice": ["1", "2", "3", "4", "5"]
	}`

	output, err := j2s.Convert(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output)

	// Output: type J2S1 struct {
	// 	CamelCase bool        `json:"CamelCase"`
	// 	KebabCase interface{} `json:"kebab-case"`
	// 	Map       J2S2        `json:"map"`
	// 	Slice     []string    `json:"slice"`
	// 	SnakeCase int         `json:"snake_case"`
	// 	Title     string      `json:"title"`
	// }
	//
	// type J2S2 struct {
	// 	Child1 string `json:"child1"`
	// 	Child2 int    `json:"child2"`
	// }
```
