/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: objutil.go
 * @time: 2017/9/30 10:38
 */
package comm_libs

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Interface2Int(v interface{}) int {
	if v == nil {
		return 0
	}
	switch d := v.(type) {
	case string:
		v, _ := strconv.Atoi(d)
		return v
	case float32, float64:
		return int(reflect.ValueOf(d).Float())
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(d).Uint())
	}
	return 0
}

func Interface2Int64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch d := v.(type) {
	case string:
		t, _ := strconv.ParseInt(d, 10, 64)
		return t
	case float32, float64:
		return int64(reflect.ValueOf(d).Float())
	case int, int8, int16, int32, int64:
		return int64(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(d).Uint())
	}
	return 0
}

func Interface2Uint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}
	switch d := v.(type) {
	case string:
		t, _ := strconv.ParseUint(d, 10, 64)
		return t
	case float32, float64:
		return uint64(reflect.ValueOf(d).Float())
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return uint64(reflect.ValueOf(d).Uint())
	}
	return 0
}

// Float64 coerces into a float64
func Interface2Float64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch d := v.(type) {
	case string:
		t, _ := strconv.ParseFloat(d, 64)
		return t
	case float32, float64:
		return reflect.ValueOf(d).Float()
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(d).Uint())
	}
	return 0
}

func Interface2Bool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch d := v.(type) {
	case bool:
		return d
	case string:
		t, _ := strconv.ParseBool(d)
		return t
	case float32, float64:
		return reflect.ValueOf(d).Float() > 0.0
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(d).Int() > 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(d).Uint() > 0
	}
	return false
}

func Interface2String(v interface{}) string {
	if v == nil {
		return ""
	}
	switch d := v.(type) {
	case string:
		return d
	default:
		return fmt.Sprintf("%v", d)
	}
	// return ""
}

// 从target中取出对应字段的值，可以是字段名，也可以是json tag
func GetFieldValueFromTarget(field string, tar interface{}) interface{} {
	myref := reflect.ValueOf(tar)
	typeOfType := myref.Type()
	if myref.Kind() == reflect.Ptr {
		myref = myref.Elem()
		typeOfType = myref.Type()
	}

	// 获取字段名，优先使用json tag, 然后使用xorm tag, 如果没有则使用字段名的小驼峰格式
	getFieldTag := func(fieldType reflect.StructField) string {
		fieldName := fieldType.Tag.Get("json")
		if len(fieldName) < 1 || fieldName == "-" {
			fieldName = fieldType.Name
			return LowCaseString(fieldName)
		}

		if trr := strings.Split(fieldName, ","); len(trr) > 1 {
			fieldName = strings.TrimSpace(trr[0]) // 过滤掉omitempty
		}
		return fieldName
	}

	if _, ok := typeOfType.FieldByName(field); ok {
		fe := myref.FieldByName(field)
		return fe.Interface()
	} else {
		for i := 0; i < myref.NumField(); i++ {
			tg := getFieldTag(typeOfType.Field(i))
			if tg == field {
				fe := myref.Field(i)
				return fe.Interface()
			}
		}
	}
	return nil
}

// 获取函数的地址
func GetFunPointer(i interface{}) uintptr {
	ref := reflect.ValueOf(i)
	if ref.Kind() == reflect.Func {
		// fmt.Println(runtime.FuncForPC(ref.Pointer()).Name(), ref.Pointer())
		return ref.Pointer()
	}
	return 0
}

func Map2KVPairs(data map[string]interface{}) []interface{} {
	out := make([]interface{}, 0, len(data)*2)
	for k, v := range data {
		out = append(out, k, v)
	}
	return out
}

func GetMapKeys(in map[string]interface{}) []string {
	out := make([]string, 0, len(in))
	for k := range in {
		out = append(out, k)
	}
	return out
}

func MapMerge(in map[string]interface{}, others ...map[string]interface{}) map[string]interface{} {
	if len(in) == 0 {
		in = map[string]interface{}{}
	}

	for _, other := range others {
		for k, v := range other {
			in[k] = v
		}
	}
	return in
}

func MapString2Int(in map[string]string) map[string]int {
	out := make(map[string]int, len(in))
	for k, v := range in {
		tmp, _ := strconv.Atoi(v)
		out[k] = tmp
	}
	return out
}

func MapInterface2Int(in map[string]interface{}) map[string]int {
	out := make(map[string]int, len(in))
	for k, v := range in {
		out[k] = Interface2Int(v)
	}
	return out
}

func MapInterface2Float(in map[string]interface{}) map[string]float64 {
	out := make(map[string]float64, len(in))
	for k, v := range in {
		out[k] = Interface2Float64(v)
	}
	return out
}

func MapInterface2String(in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = Interface2String(v)
	}
	return out
}

func MapInterface2Bool(in map[string]interface{}) map[string]bool {
	out := make(map[string]bool, len(in))
	for k, v := range in {
		out[k] = Interface2Bool(v)
	}
	return out
}

func MapString2Interface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
