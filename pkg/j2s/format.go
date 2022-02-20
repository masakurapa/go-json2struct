package j2s

import (
	"encoding/json"
	"fmt"
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
		return "", fmt.Errorf("json unmarshal Error: %v", err)
	}

	switch v := val.(type) {
	case bool:
		return fmt.Sprintf(typeFormat, 1, boolType), nil
	case float64:
		return fmt.Sprintf(typeFormat, 1, numberType(v)), nil
	case nil:
		return fmt.Sprintf(typeFormat, 1, interfaceType), nil
	case map[string]interface{}:
		return formatStruct(v)
	default:
		// it shouldn't be possible, but just in case
		return "", fmt.Errorf("unsupported %v", v)
	}
}

func numberType(v float64) string {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		return intType
	}
	return floatType
}
