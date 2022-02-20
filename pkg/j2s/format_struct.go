package j2s

import (
	"fmt"
	"go/format"
	"strings"
)

func formatStruct(v map[string]interface{}) (string, error) {
	f := structFormatter{}
	f.appendStruct(v)

	structs := make([]string, 0, len(f.structs))
	for i := len(f.structs) - 1; i >= 0; i-- {
		st := f.structs[i]
		structs = append(structs, fmt.Sprintf(typeFormat, st.no, st.string()))
	}

	code := strings.TrimSpace(strings.Join(structs, "\n\n"))
	b, err := format.Source([]byte(code))
	if err != nil {
		return "", fmt.Errorf("code format error: %v", err)
	}
	return strings.TrimSpace(string(b)), nil
}

type structFormatter struct {
	cnt     int
	structs []structType
}

func (f *structFormatter) nextNo() int {
	return f.cnt + 1
}

func (f *structFormatter) nextStruct() structType {
	f.cnt++
	return structType{no: f.cnt}
}

func (f *structFormatter) appendStruct(v map[string]interface{}) {
	st := f.nextStruct()
	for key, val := range v {
		st.fields = append(st.fields, fmt.Sprintf(structFieldFormat, f.field(key), f.getType(val), key))
	}
	f.structs = append(f.structs, st)
}

func (f *structFormatter) field(s string) string {
	return link.ReplaceAllStringFunc(s, func(s string) string {
		ss := strings.Replace(strings.Replace(s, "_", "", -1), "-", "", -1)
		return strings.ToUpper(ss)
	})
}

func (f *structFormatter) getType(v interface{}) string {
	switch i := v.(type) {
	case bool:
		return boolType
	case string:
		return stringType
	case float64:
		return numberType(i)
	case map[string]interface{}:
		t := fmt.Sprintf(typeName, f.nextNo())
		f.appendStruct(i)
		return t
	case []interface{}:
		return f.getSliceType(i)
	}
	return interfaceType
}

func (f *structFormatter) getSliceType(v []interface{}) string {
	ret := ""
	for _, vv := range v {
		t := f.getType(vv)

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

type structType struct {
	no     int
	fields []string
}

func (st *structType) string() string {
	return fmt.Sprintf("struct {\n%s}", strings.Join(st.fields, "\n"))
}
