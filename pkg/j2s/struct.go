package j2s

import (
	"fmt"
	"strings"
)

func (c *converter) getStructType(no int, v map[string]interface{}) string {
	fields := make([]string, 0, len(v))
	for key, val := range v {
		fields = append(fields, fmt.Sprintf(structFieldFormat, c.structField(key), c.getType(no, val), key))
	}
	return "struct {\n" + strings.Join(fields, "\n") + "}"
}

func (c *converter) structField(s string) string {
	return link.ReplaceAllStringFunc(s, func(s string) string {
		ss := strings.Replace(strings.Replace(s, "_", "", -1), "-", "", -1)
		return strings.ToUpper(ss)
	})
}

func (c *converter) getSliceType(no int, v []interface{}) string {
	ret := ""
	for _, vv := range v {
		t := c.getType(no, vv)

		if ret == "" {
			ret = t
			continue
		}

		if ret != t {
			ret = interfaceType
			break
		}
	}
	return "[]" + ret
}
