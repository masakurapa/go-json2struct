package j2s

import (
	"encoding/json"
	"errors"
	"go/format"
	"regexp"
	"strconv"
	"strings"
)

const (
	typeName          = "J2S%d"
	typeFormat        = "type " + typeName + " %s"
	structFieldFormat = "%s %s `json:\"%s\"`\n"

	boolType      = "bool"
	intType       = "int"
	floatType     = "float64"
	stringType    = "string"
	interfaceType = "interface{}"
)

var (
	link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])|-([A-Za-z])")
)

func Format(s string) (string, error) {
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
	c.appendTypes(1, v)
	return c.toString()
}

func (c *converter) appendTypes(no int, v interface{}) {
	typeName := c.getType(no, v)
	code := "type J2S" + strconv.Itoa(no) + " " + typeName
	c.types = append(c.types, typeInfo{
		no:   no,
		code: code,
	})
}

func (c *converter) getType(no int, v interface{}) string {
	switch vv := v.(type) {
	case bool:
		return boolType
	case string:
		return stringType
	case float64:
		return c.getNumberTyp(vv)
	case map[string]interface{}:
		return c.getStructType(no, vv)
	case []interface{}:
		return c.getSliceType(no, vv)
	}
	return interfaceType
}

func (c *converter) getNumberTyp(v float64) string {
	// TODO: there has to be a better way
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		return intType
	}
	return floatType
}

func (c *converter) toString() (string, error) {
	codes := make([]string, len(c.types))
	for _, ti := range c.types {
		codes[ti.no-1] = ti.code
	}

	code := strings.TrimSpace(strings.Join(codes, "\n\n"))
	b, err := format.Source([]byte(code))
	if err != nil {
		return "", errors.New("code format error: " + err.Error())
	}
	return strings.TrimSpace(string(b)), nil
}
