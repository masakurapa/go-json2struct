package j2s_test

import (
	"fmt"
	"testing"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func TestConvert(t *testing.T) {
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
		// FIXME: "test5": need to fix the bug with the int type.
		{
			name: "struct",
			s:    `{"test1":null,"test2":"1","test3":2,"test4":123.456,"test5":123.000,"test6":true}`,
			expected: `type J2S1 struct {
	Test1 interface{} ` + "`json:\"test1\"`" + `
	Test2 string      ` + "`json:\"test2\"`" + `
	Test3 int         ` + "`json:\"test3\"`" + `
	Test4 float64     ` + "`json:\"test4\"`" + `
	Test5 int         ` + "`json:\"test5\"`" + `
	Test6 bool        ` + "`json:\"test6\"`" + `
}`,
		},
		{
			name: "struct - Snake-Case key",
			s:    `{"test_key_a":null}`,
			expected: `type J2S1 struct {
	TestKeyA interface{} ` + "`json:\"test_key_a\"`" + `
}`,
		},
		{
			name: "struct - Kebab-Case key",
			s:    `{"test-key-a":null}`,
			expected: `type J2S1 struct {
	TestKeyA interface{} ` + "`json:\"test-key-a\"`" + `
}`,
		},
		{
			name: "struct - number separate key",
			s:    `{"test1key2a":null}`,
			expected: `type J2S1 struct {
	Test1key2a interface{} ` + "`json:\"test1key2a\"`" + `
}`,
		},

		{
			name: "struct - slice value",
			s:    `{"test":["1", "2", "3"]}`,
			expected: `type J2S1 struct {
	Test []string ` + "`json:\"test\"`" + `
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
		{
			name: "struct - slice on map value",
			s:    `{"test":[{"fuga":"12345"}]}`,
			expected: `type J2S1 struct {
	Test []J2S2 ` + "`json:\"test\"`" + `
}

type J2S2 struct {
	Fuga string ` + "`json:\"fuga\"`" + `
}`,
		},

		{
			name:     "slice - empty",
			s:        `[]`,
			expected: `type J2S1 []interface{}`,
		},
		{
			name:     "slice - null value",
			s:        `[null, null, null]`,
			expected: `type J2S1 []interface{}`,
		},
		{
			name:     "slice - string value",
			s:        `["1", "2", "3"]`,
			expected: `type J2S1 []string`,
		},
		{
			name:     "slice - int value",
			s:        `[1, 2, 3]`,
			expected: `type J2S1 []int`,
		},
		{
			name:     "slice - float value",
			s:        `[1.1, 2.2, 3.3]`,
			expected: `type J2S1 []float64`,
		},
		{
			name:     "slice - slice value",
			s:        `[["1","2"]]`,
			expected: `type J2S1 [][]string`,
		},
		{
			name:     "slice - slice on slice value",
			s:        `[[["1","2"]]]`,
			expected: `type J2S1 [][][]string`,
		},
		{
			name: "slice - map value",
			s:    `[{"test":"1"}]`,
			expected: `type J2S1 []J2S2

type J2S2 struct {
	Test string ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "slice - multiple map value1",
			s:    `[{"test":"1"},{"test":"2","food":"apple"},{"test":"2","drink":"beer"}]`,
			expected: `type J2S1 []J2S2

type J2S2 struct {
	Drink *string ` + "`json:\"drink\"`" + `
	Food  *string ` + "`json:\"food\"`" + `
	Test  string  ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "slice - multiple map value2",
			s:    `[{"test":"1","food":""},{"test":2,"food":"apple"},{"food":null}]`,
			expected: `type J2S1 []J2S2

type J2S2 struct {
	Food interface{}  ` + "`json:\"food\"`" + `
	Test *interface{} ` + "`json:\"test\"`" + `
}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := j2s.Convert(tc.s)

			if err == nil && tc.err == nil {
				if actual != tc.expected {
					t.Errorf("Convert() returns: \n%v\nwant: \n%v", actual, tc.expected)
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

func TestConvertWithOption(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		opt      j2s.Option
		expected string
		err      error
	}{
		{
			name: "UseTag is false",
			s:    `{"test":"1"}`,
			opt:  j2s.Option{UseTag: false, TagName: "example"},
			expected: `type J2S1 struct {
	Test string
}`,
		},
		{
			name: "TagName is empty string",
			s:    `{"test":"1"}`,
			opt:  j2s.Option{UseTag: true, TagName: ""},
			expected: `type J2S1 struct {
	Test string ` + "`json:\"test\"`" + `
}`,
		},
		{
			name: "specify non-json for TagName",
			s:    `{"test":"1"}`,
			opt:  j2s.Option{UseTag: true, TagName: "bson"},
			expected: `type J2S1 struct {
	Test string ` + "`bson:\"test\"`" + `
}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := j2s.ConvertWithOption(tc.s, tc.opt)

			if err == nil && tc.err == nil {
				if actual != tc.expected {
					t.Errorf("Convert() returns: \n%v\nwant: \n%v", actual, tc.expected)
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
