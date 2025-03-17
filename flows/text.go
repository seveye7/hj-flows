package flows

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"hj-flows/utils"
)

func UnMarshalBytes[T any](arr [][]byte) ([]*T, []string, error) {
	strs := []string{}
	for _, v := range arr {
		strs = append(strs, utils.B2s(v))
	}
	return UnMarshal[T](strs)
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
		for i := 0; i < elem.NumField(); i++ {
			if i >= len(fieldArrs) {
				break
			}
			// 判断类型 float64, int64, uint64, string
			// 如果类型是 string, 则直接赋值
			kind := elem.Field(i).Type().Kind()
			switch kind {
			case reflect.Float64, reflect.Float32:
				if fieldArrs[i] == "" {
					elem.Field(i).SetFloat(0)
					break
				}
				f, err := strconv.ParseFloat(fieldArrs[i], 64)
				if err != nil {
					ok = false
					break
					// return reflect.Value{}, fmt.Errorf("parse float64 error: %v", err)
				}
				elem.Field(i).SetFloat(f)
			case reflect.Int64, reflect.Int32:
				if fieldArrs[i] == "" {
					elem.Field(i).SetInt(0)
					break
				}
				f, err := strconv.ParseInt(fieldArrs[i], 10, 64)
				if err != nil {
					ok = false
					break
					// return reflect.Value{}, fmt.Errorf("parse int64 error: %v", err)
				}
				elem.Field(i).SetInt(f)
			case reflect.Uint64, reflect.Uint32:
				if fieldArrs[i] == "" {
					elem.Field(i).SetUint(0)
					break
				}

				f, err := strconv.ParseUint(fieldArrs[i], 10, 64)
				if err != nil {
					ok = false
					break
					// return reflect.Value{}, fmt.Errorf("parse uint64 error: %v", err)
				}
				elem.Field(i).SetUint(f)
			case reflect.String:
				elem.Field(i).SetString(fieldArrs[i])
			}
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
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var arr []string
	for i := 0; i < v.NumField(); i++ {
		arr = append(arr, fmt.Sprintf("%v", v.Field(i).Interface()))
	}
	return strings.Join(arr, ";")
}

func MarshalBytes(i any) [][]byte {
	strs := Marshal(i)
	ret := make([][]byte, len(strs))
	for i, v := range strs {
		ret[i] = utils.S2b(v)
	}
	return ret
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
