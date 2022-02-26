package j2s

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/format"
	"regexp"
	"strconv"
	"strings"
)

const (
	firstNo = 1
)

var (
	link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])|-([A-Za-z])")
)

// Convert returns a string that converts a JSON string into a Go structure.
//
// return an error if the string is invalid as JSON.
func Convert(s string) (string, error) {
	var val interface{}
	if err := json.Unmarshal([]byte(s), &val); err != nil {
		// if you only have a string like "hoge", you will get an error here.
		return "", errors.New("json unmarshal Error: " + err.Error())
	}

	conv := converter{}
	return conv.toStruct(val)
}

type converter struct {
	types []typeInfo
}

type typeInfo struct {
	no   int
	code string
}

func (c *converter) toStruct(v interface{}) (string, error) {
	c.appendTypes(firstNo, v)
	return c.toString()
}

func (c *converter) append(no int, code string) {
	c.types = append(c.types, typeInfo{
		no:   no,
		code: code,
	})
}

func (c *converter) appendTypes(no int, v interface{}) {
	typeName := c.getType(no, v)
	code := "type J2S" + strconv.Itoa(no) + " " + typeName
	c.append(no, code)
}

func (c *converter) getType(no int, v interface{}) string {
	switch vv := v.(type) {
	case bool:
		return "bool"
	case string:
		return "string"
	case float64:
		return c.getNumberTyp(vv)
	case map[string]interface{}:
		return c.getStructType(no, vv)
	case []interface{}:
		return c.getSliceType(no, vv)
	}
	return "interface{}"
}

func (c *converter) getNumberTyp(v float64) string {
	// TODO: there has to be a better way
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		return "int"
	}
	return "float64"
}

func (c *converter) getStructType(no int, v map[string]interface{}) string {
	buf := bytes.Buffer{}
	buf.WriteString("struct {\n")

	for key, val := range v {
		nextNo := no + 1
		typeStr := c.getType(nextNo, val)
		if strings.HasPrefix(typeStr, "struct {") {
			structName := "J2S" + strconv.Itoa(nextNo)
			c.append(nextNo, "type "+structName+" "+typeStr)
			typeStr = structName
		}
		buf.WriteString(c.structField(key) + " " + typeStr)
		buf.WriteString(" `json:\"" + key + "\"`")
		buf.WriteString("\n")
	}

	buf.WriteString("}")
	return buf.String()
}

func (c *converter) structField(s string) string {
	return link.ReplaceAllStringFunc(s, func(s string) string {
		ss := strings.Replace(strings.Replace(s, "_", "", -1), "-", "", -1)
		return strings.ToUpper(ss)
	})
}

func (c *converter) getSliceType(no int, v []interface{}) string {
	if len(v) == 0 {
		return "[]interface{}"
	}

	ret := ""
	for _, vv := range v {
		t := "interface{}"
		switch vvv := vv.(type) {
		case bool:
			t = "bool"
		case string:
			t = "string"
		case float64:
			t = c.getNumberTyp(vvv)
		case map[string]interface{}:
			t = "map[string]interface{}"
		case []interface{}:
			t = "[]interface{}"
		}

		if ret == "" {
			ret = t
			continue
		}

		if ret != t {
			ret = "interface{}"
			break
		}
	}
	return "[]" + ret
}

func (c *converter) toString() (string, error) {
	codes := make([]string, len(c.types))
	for _, ti := range c.types {
		codes[ti.no-1] = ti.code
	}

	code := strings.Join(codes, "\n\n")
	b, err := format.Source([]byte(code))
	if err != nil {
		return "", errors.New("code format error: " + err.Error())
	}
	return strings.TrimSpace(string(b)), nil
}
