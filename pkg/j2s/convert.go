package j2s

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/format"
	"regexp"
	"sort"
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
	no int
	// string of type
	// not used for struct type
	typeStr string

	isStruct     bool
	structFields []structField
}

type structField struct {
	name    string
	typeStr string
}

func (c *converter) toStruct(v interface{}) (string, error) {
	c.append(c.getTypeInfo(firstNo, v))
	return c.toString()
}

func (c *converter) append(v typeInfo) {
	c.types = append(c.types, v)
}

func (c *converter) getTypeInfo(no int, v interface{}) typeInfo {
	ti := typeInfo{no: no}
	switch vv := v.(type) {
	case bool:
		ti.typeStr = "bool"
	case string:
		ti.typeStr = "string"
	case float64:
		ti.typeStr = c.getNumberTyp(vv)
	case map[string]interface{}:
		// overwrite "ti"
		ti = c.getStructTypeInfo(no, vv)
	case []interface{}:
		ti.typeStr = c.getSliceType(no, vv)
	default:
		ti.typeStr = "interface{}"
	}
	return ti
}

func (c *converter) getNumberTyp(v float64) string {
	// TODO: there has to be a better way
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		return "int"
	}
	return "float64"
}

func (c *converter) getStructTypeInfo(no int, v map[string]interface{}) typeInfo {
	ti := typeInfo{
		no:           no,
		isStruct:     true,
		structFields: make([]structField, 0, len(v)),
	}

	for key, vv := range v {
		nextNo := no + 1
		fieldInfo := c.getTypeInfo(nextNo, vv)
		typeStr := fieldInfo.typeStr
		if fieldInfo.isStruct {
			structName := "J2S" + strconv.Itoa(nextNo)
			c.append(fieldInfo)
			typeStr = structName
		}

		ti.structFields = append(ti.structFields, structField{
			name:    key,
			typeStr: typeStr,
		})
	}

	return ti
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
			nextNo := no
			if no == firstNo {
				nextNo++
			}

			c.append(c.getTypeInfo(nextNo, vvv))
			t = "J2S" + strconv.Itoa(nextNo)
		case []interface{}:
			t = c.getSliceType(no+1, vvv)
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
		code := "type J2S" + strconv.Itoa(ti.no) + " "
		if ti.isStruct {
			code += c.toStructString(ti)
		} else {
			code += ti.typeStr
		}
		codes[ti.no-1] = code
	}

	code := strings.Join(codes, "\n\n")
	b, err := format.Source([]byte(code))
	if err != nil {
		return "", errors.New("code format error: " + err.Error())
	}
	return strings.TrimSpace(string(b)), nil
}

func (c *converter) toStructString(ti typeInfo) string {
	// sort by key name in asc
	keys := make([]string, 0, len(ti.structFields))
	fields := make(map[string]structField)
	for _, field := range ti.structFields {
		fields[field.name] = field
		keys = append(keys, field.name)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	buf := bytes.Buffer{}
	buf.WriteString("struct {\n")

	for _, key := range keys {
		field := fields[key]
		buf.WriteString(c.structField(field.name) + " " + field.typeStr)
		buf.WriteString(" `json:\"" + field.name + "\"`")
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
