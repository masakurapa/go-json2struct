package j2s_test

import (
	"fmt"
	"testing"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func TestFormat(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		expected string
		err      error
	}{
		{
			name:     "null",
			s:        "null",
			expected: "type J2S1 interface{}",
		},
		{
			name:     "string",
			s:        "aaaaa",
			expected: "",
			err:      fmt.Errorf("json unmarshal Error: invalid character 'a' looking for beginning of value"),
		},
		{
			name:     "int",
			s:        "12345",
			expected: "type J2S1 int",
		},
		{
			name:     "float64",
			s:        "12.345",
			expected: "type J2S1 float64",
		},
		{
			name:     "bool",
			s:        "true",
			expected: "type J2S1 bool",
		},

		{
			name: "struct - empty",
			s:    `{}`,
			expected: `type J2S1 struct {
}`,
		},
		{
			name: "struct - null value",
			s:    `{"test":null}`,
			expected: `type J2S1 struct {
	Test interface{} ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - null value(snake)",
			s:    `{"test_key_a":null}`,
			expected: `type J2S1 struct {
	TestKeyA interface{} ` + "`json:\"test_key_a\"`" + `
}`,
		},
		{
			name: "struct - null value(kebab)",
			s:    `{"test-key-a":null}`,
			expected: `type J2S1 struct {
	TestKeyA interface{} ` + "`json:\"test-key-a\"`" + `
}`,
		},
		{
			name: "struct - null value(number separate)",
			s:    `{"test1key2a":null}`,
			expected: `type J2S1 struct {
	Test1key2a interface{} ` + "`json:\"test1key2a\"`" + `
}`,
		},
		{
			name: "struct - string value",
			s:    `{"test":"1"}`,
			expected: `type J2S1 struct {
	Test string ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - int value",
			s:    `{"test":1}`,
			expected: `type J2S1 struct {
	Test int ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - float value",
			s:    `{"test":123.456}`,
			expected: `type J2S1 struct {
	Test float64 ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - bool value",
			s:    `{"test":true}`,
			expected: `type J2S1 struct {
	Test bool ` + "`json:\"test\"`" + `
}`,
		},

		{
			name: "struct - null slice value",
			s:    `{"test":[null, null, null]}`,
			expected: `type J2S1 struct {
	Test []interface{} ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - string slice value",
			s:    `{"test":["1", "2", "3"]}`,
			expected: `type J2S1 struct {
	Test []string ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - int slice value",
			s:    `{"test":[1, 2, 3]}`,
			expected: `type J2S1 struct {
	Test []int ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - float slice value",
			s:    `{"test":[1.1, 2.2, 3.3]}`,
			expected: `type J2S1 struct {
	Test []float64 ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - bool slice value",
			s:    `{"test":[true, false, true]}`,
			expected: `type J2S1 struct {
	Test []bool ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "struct - multilpe value type",
			s:    `{"test":[null, "1", 2, 3.3, true]}`,
			expected: `type J2S1 struct {
	Test []interface{} ` + "`json:\"test\"`" + `
}`,
		},

		{
			name: "struct - map value",
			s:    `{"test":{"hoge":"fuga"}}`,
			expected: `type J2S1 struct {
	Test J2S2 ` + "`json:\"test\"`" + `
}

type J2S2 struct {
	Hoge string ` + "`json:\"hoge\"`" + `
}`,
		},
		{
			name: "struct - nested map value",
			s:    `{"test":{"hoge":{"fuga":"12345"}}}`,
			expected: `type J2S1 struct {
	Test J2S2 ` + "`json:\"test\"`" + `
}

type J2S2 struct {
	Hoge J2S3 ` + "`json:\"hoge\"`" + `
}

type J2S3 struct {
	Fuga string ` + "`json:\"fuga\"`" + `
}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := j2s.Format(tc.s)

			if err == nil && tc.err == nil {
				if actual != tc.expected {
					t.Errorf("Format() returns: \n%v\nwant: \n%v", actual, tc.expected)
				}
				return
			}

			if err != nil && tc.err != nil {
				if err.Error() != tc.err.Error() {
					t.Fatalf("error returns %v, want %v", err, tc.err)
				}
				return
			}

			t.Fatalf("error returns %v, want %v", err, tc.err)
		})
	}
}
