package flows

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"hj-flows/utils"
)

func MarshalBytes(i any) [][]byte {
	strs := Marshal(i)
	ret := make([][]byte, len(strs))
	for i, v := range strs {
		ret[i] = utils.S2b(v)
	}
	return ret
}

func UnMarshalBytes[T any](arr [][]byte) ([]*T, []string, error) {
	strs := []string{}
	for _, v := range arr {
		strs = append(strs, utils.B2s(v))
	}
	return UnMarshal[T](strs)
}

func parseElem(elem reflect.Value, i int, fieldStr string) bool {
	field := elem.Field(i)
	kind := field.Type().Kind()
	isPtr := kind == reflect.Ptr
	if isPtr {
		// 空不赋值，*string也是空指针
		if fieldStr == "" || fieldStr == "nil" {
			return true
		}
		// new
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		kind = field.Elem().Type().Kind()
		field = field.Elem()
	}
	switch kind {
	case reflect.Float64, reflect.Float32:
		if fieldStr == "" {
			elem.Field(i).SetFloat(0)
			break
		}
		f, err := strconv.ParseFloat(fieldStr, 64)
		if err != nil {
			return false
		}
		field.SetFloat(f)
	case reflect.Int64, reflect.Int32:
		if fieldStr == "" {
			field.SetInt(0)
		}
		f, err := strconv.ParseInt(fieldStr, 10, 64)
		if err != nil {
			return false
		}

		field.SetInt(f)
	case reflect.Uint64, reflect.Uint32:
		if fieldStr == "" {
			field.SetUint(0)
		}

		f, err := strconv.ParseUint(fieldStr, 10, 64)
		if err != nil {
			return false
			// return reflect.Value{}, fmt.Errorf("parse uint64 error: %v", err)
		}
		field.SetUint(f)
	case reflect.String:
		field.SetString(fieldStr)
	}
	return true
}

func UnMarshal[T any](arr []string) ([]*T, []string, error) {
	slice := []*T{}
	dirty := []string{}

	// r := reflect.TypeOf(t)
	for _, str := range arr {
		t := new(T)
		elem := reflect.ValueOf(t)
		if reflect.TypeOf(t).Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		ok := true
		fieldArrs := strings.Split(str, ";")
		start := 0
		n := 1
		if elem.Type().Field(0).Anonymous {
			for i := 0; i < elem.Field(0).NumField(); i++ {
				if i >= len(fieldArrs) {
					break
				}
				if !parseElem(elem.Field(0), i, fieldArrs[i]) {
					ok = false
					break
				}
			}
			start = 1
			n = elem.Field(0).NumField()
		}
		if !ok {
			dirty = append(dirty, str)
			continue
		}
		for i := start; i < elem.NumField(); i++ {
			if i+n > len(fieldArrs) {
				break
			}
			// 判断类型 float64, int64, uint64, string
			// 如果类型是 string, 则直接赋值
			parseElem(elem, i, fieldArrs[i+n-1])
		}

		if ok {
			slice = append(slice, t)
		} else {
			dirty = append(dirty, str)
		}
	}

	return slice, dirty, nil
}

func StructToString(s interface{}) string {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	var arr []string
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Anonymous {
			arr = append(arr, StructToString(v.Field(i).Interface()))
			continue
		} else if v.Field(i).Kind() == reflect.Ptr {
			if v.Field(i).IsNil() {
				// 指针类型
				arr = append(arr, "")
			} else {
				arr = append(arr, fmt.Sprint(v.Field(i).Elem().Interface()))
			}
			continue
		}
		arr = append(arr, fmt.Sprintf("%v", v.Field(i).Interface()))
	}
	return strings.Join(arr, ";")
}

func Marshal(i any) []string {
	var arr []string
	if i == nil {
		return arr
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array:
		slice := reflect.ValueOf(i)
		for i := 0; i < slice.Len(); i++ {
			arr = append(arr, StructToString(slice.Index(i).Interface()))
		}
	default:
		arr = append(arr, StructToString(i))
	}
	return arr
}

func StructToValue(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var arr []string
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Type().Kind() {
		case reflect.String:
			arr = append(arr, "'"+v.Field(i).Interface().(string)+"'")
		default:
			arr = append(arr, fmt.Sprintf("%v", v.Field(i).Interface()))
		}
	}
	return "(" + strings.Join(arr, ",") + ")"
}

func MarshalValues(i any) []string {
	var arr []string
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array:
		slice := reflect.ValueOf(i)
		for i := 0; i < slice.Len(); i++ {
			arr = append(arr, StructToValue(slice.Index(i).Interface()))
		}
	default:
		arr = append(arr, StructToValue(i))
	}
	return arr
}

type Model interface {
	TableName() string
	PartitionKey() string
}

const (
	sqlInsertStr = "INSERT INTO "
	sqlValueStr  = ") VALUES "
)

// StructToInsertSql 生成批量插入sql，注意指针必须都设置
func StructToInsertSql(dbName string, models []Model) string {
	if len(models) == 0 {
		return ""
	}

	var valuess, fields []string
	for i0, model := range models {
		v := reflect.ValueOf(model)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		var values []string
		for i := 0; i < v.NumField(); i++ {
			dbTag := v.Type().Field(i).Tag.Get("db")
			if dbTag == "" || dbTag == "-" {
				continue
			}
			if i0 == 0 {
				fields = append(fields, "`"+dbTag+"`")
			}

			// 判断是否为空
			if v.Field(i).Type().Kind() == reflect.Ptr && v.Field(i).IsNil() {
				values = append(values, "NULL")
				continue
			}

			//
			elem := v.Field(i)
			if v.Field(i).Type().Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			switch elem.Type().Kind() {
			case reflect.String:
				values = append(values, "'"+elem.Interface().(string)+"'")
			default:
				values = append(values, fmt.Sprintf("%v", elem.Interface()))
			}
		}
		valuess = append(valuess, "("+strings.Join(values, ",")+")")
	}
	filedStr := strings.Join(fields, ",")
	valuesStr := strings.Join(valuess, ",")

	var buff strings.Builder
	// 计算sql长度
	size := len(sqlInsertStr) + len(models[0].TableName()) + 1 + len(filedStr) + len(sqlValueStr) + len(valuesStr) + 1
	if dbName != "" {
		size += len(dbName) + 1
	}
	buff.Grow(size)

	// 拼接sql
	buff.WriteString(sqlInsertStr)
	if dbName != "" {
		buff.WriteString(dbName)
		buff.WriteString(".")
	}
	buff.WriteString(models[0].TableName())
	buff.WriteString("(")
	buff.WriteString(filedStr)
	buff.WriteString(sqlValueStr)
	buff.WriteString(valuesStr)
	buff.WriteString(";")
	return buff.String()
}
