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

	boolType         = "bool"
	stringType       = "string"
	intType          = "int"
	floatType        = "float64"
	interfaceType    = "interface{}"
	structNamePrefix = "J2S"

	defaultTag = "json"
	omitempty  = ",omitempty"
)

// Omitempty is an optional type that outputs "omitempty"
type Omitempty int

const (
	// OmitemptyNone does not always output "omitempty"
	OmitemptyNone Omitempty = iota
	// OmitemptyForAll always output "omitempty"
	OmitemptyForAll
	// OmitemptyPtrOnly outputs "omitempty" only if it is a pointer type
	OmitemptyPtrOnly
)

var (
	link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])|-([A-Za-z])")
)

// Convert uses the default options and
// returns a string of JSON strings converted to Go structures.
//
// The default values are as follows.
//     Option {
//         UseTag:    true,
//         TagName:   "json",
//         Omitempty: OmitemptyNone,
//     }
//
// return an error if the string is invalid as JSON.
func Convert(s string) (string, error) {
	return ConvertWithOption(s, Option{
		UseTag:    true,
		TagName:   defaultTag,
		Omitempty: OmitemptyNone,
	})
}

// ConvertWithOption uses the specified options and
// returns a string that is a JSON string converted to a Go structure.
//
// return an error if the string is invalid as JSON.
func ConvertWithOption(s string, opt Option) (string, error) {
	var val interface{}
	if err := json.Unmarshal([]byte(s), &val); err != nil {
		// if you only have a string like "hoge", you will get an error here.
		return "", errors.New("json unmarshal Error: " + err.Error())
	}

	conv := converter{opt: opt}
	return conv.toStruct(val)
}

// Option is an option to customize output results
type Option struct {
	// UseTag outputs tag if true
	UseTag bool

	// TagName is a tag name. (default is "json")
	//
	// if empty, the default value is used.
	//
	// if "UseTag" is false, the value is not used.
	TagName string

	// Omitempty is an optional type that outputs "omitempty". (default is non output)
	Omitempty Omitempty
}

type converter struct {
	types []typeInfo
	opt   Option
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
	isPtr   bool
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
		ti.typeStr = boolType
	case string:
		ti.typeStr = stringType
	case float64:
		ti.typeStr = c.getNumberTyp(vv)
	case map[string]interface{}:
		// overwrite "ti"
		ti = c.getStructTypeInfo(no, vv)
	case []interface{}:
		ti.typeStr = c.getSliceType(no, vv)
	default:
		ti.typeStr = interfaceType
	}
	return ti
}

func (c *converter) getNumberTyp(v float64) string {
	// TODO: there has to be a better way
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		return intType
	}
	return floatType
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
			structName := structNamePrefix + strconv.Itoa(nextNo)
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
		return "[]" + interfaceType
	}

	ret := ""
	for _, vv := range v {
		t := interfaceType
		switch vvv := vv.(type) {
		case bool:
			t = boolType
		case string:
			t = stringType
		case float64:
			t = c.getNumberTyp(vvv)
		case map[string]interface{}:
			nextNo := no
			if no == firstNo {
				nextNo++
			}

			c.append(c.getTypeInfo(nextNo, vvv))
			t = structNamePrefix + strconv.Itoa(nextNo)
		case []interface{}:
			t = c.getSliceType(no+1, vvv)
		}

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

func (c *converter) toString() (string, error) {
	// grouping by same number
	groups := make(map[int][]typeInfo)
	for _, ti := range c.types {
		groups[ti.no] = append(groups[ti.no], ti)
	}

	codes := make([]string, len(c.types))
	for no, tis := range groups {
		code := "type " + structNamePrefix + strconv.Itoa(no) + " "

		var ti typeInfo
		if len(tis) == 1 {
			ti = tis[0]
		} else {
			// The only case that gets us here is a
			// slice of a structure like `[{"test":"1"}, {"test":"2"}]`
			ti = c.margeTypeInfo(no, tis)
		}

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
		buf.WriteString(c.structField(field.name) + " ")

		if field.isPtr {
			buf.WriteString("*")
		}
		buf.WriteString(field.typeStr)

		buf.WriteString(c.structTag(field))
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

func (c *converter) structTag(field structField) string {
	if !c.opt.UseTag {
		return ""
	}

	tag := c.opt.TagName
	if tag == "" {
		tag = defaultTag
	}

	ret := "`" + tag + `:"` + field.name

	switch c.opt.Omitempty {
	case OmitemptyForAll:
		ret += omitempty
	case OmitemptyPtrOnly:
		if field.isPtr {
			ret += omitempty
		}
	}

	return ret + "\"`"
}

func (c *converter) margeTypeInfo(no int, tis []typeInfo) typeInfo {
	type tmpFieldInfo struct {
		cnt         int
		structField structField
	}
	appended := make(map[string]tmpFieldInfo)

	for _, ti := range tis {
		for _, field := range ti.structFields {
			tmp, ok := appended[field.name]
			if !ok {
				tmp = tmpFieldInfo{cnt: 1, structField: field}
				appended[field.name] = tmp
				continue
			}

			tmp.cnt++
			if tmp.structField.typeStr != field.typeStr {
				tmp.structField.typeStr = interfaceType
			}
			appended[field.name] = tmp
		}
	}

	// The element of slice is supposed to be a structure, so we'll set "isStruct" to true
	ret := typeInfo{
		no:           no,
		isStruct:     true,
		structFields: make([]structField, 0, len(appended)),
	}
	cnt := len(tis)

	for _, tmp := range appended {
		field := tmp.structField
		// If there are N elements, it will be false only if the field appears in all elements.
		field.isPtr = cnt != tmp.cnt
		ret.structFields = append(ret.structFields, field)
	}

	return ret
}
