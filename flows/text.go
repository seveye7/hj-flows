package flows

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"

	"hj-flows/utils"
)

func StructToBytes(s any) []byte {
	var buffer bytes.Buffer
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	arr := make([]string, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.String {
			arr = append(arr, v.String())
		} else {
			arr = append(arr, fmt.Sprintf("%v", v.Field(i).Interface()))
		}
	}
	n := 2
	for _, v := range arr {
		n += len(v) + 2
	}
	buffer.Grow(n)
	arrLen := uint16(len(arr))
	buffer.Write([]byte{byte(arrLen & 0xFF), byte(arrLen >> 8 & 0xFF)})
	// binary.Write(&buffer, binary.BigEndian, uint16(len(arr)))
	for _, v := range arr {
		vLen := uint16(len(v))
		buffer.Write([]byte{byte(vLen & 0xFF), byte(vLen >> 8 & 0xFF)})
		// binary.Write(&buffer, binary.BigEndian, uint16(len(v)))
		if len(v) > 0 {
			buffer.WriteString(v)
		}
	}
	return buffer.Bytes()
}

func MarshalBytes(i any) [][]byte {
	var arr [][]byte
	if i == nil {
		return arr
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array:
		slice := reflect.ValueOf(i)
		if slice.Len() == 0 {
			return arr
		}
		arr = make([][]byte, 0, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			arr = append(arr, StructToBytes(slice.Index(i).Interface()))
		}
	default:
		arr = make([][]byte, 0, 1)
		arr = append(arr, StructToBytes(i))
	}
	return arr
}

// splitBytes 将字节数组, 按照2字节长度+数据的方式, 分割成多个字节数组
func splitBytes(data []byte) [][]byte {
	n := binary.BigEndian.Uint16(data)
	ret := make([][]byte, 0, n)
	index := 2
	for i := 0; i < int(n); i++ {
		dataLen := binary.BigEndian.Uint16(data[i:])
		ret = append(ret, data[index+2:index+2+int(dataLen)])
		i += 2 + int(dataLen)
	}
	return ret
}

func UnMarshalBytes[T any](arr [][]byte) ([]*T, [][]byte, error) {
	slice := []*T{}
	dirty := [][]byte{}

	// r := reflect.TypeOf(t)
	for _, data := range arr {
		t := new(T)
		elem := reflect.ValueOf(t)
		if reflect.TypeOf(t).Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		ok := true
		fieldArrs := splitBytes(data)
		for i := 0; i < elem.NumField(); i++ {
			if i >= len(fieldArrs) {
				break
			}
			s := utils.B2s(fieldArrs[i])
			// 判断类型 float64, int64, uint64, string
			// 如果类型是 string, 则直接赋值
			kind := elem.Field(i).Type().Kind()
			switch kind {
			case reflect.Float64, reflect.Float32:
				if s == "" {
					elem.Field(i).SetFloat(0)
					break
				}
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					ok = false
					break
					// return reflect.Value{}, fmt.Errorf("parse float64 error: %v", err)
				}
				elem.Field(i).SetFloat(f)
			case reflect.Int64, reflect.Int32:
				if s == "" {
					elem.Field(i).SetInt(0)
					break
				}
				f, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					ok = false
					break
					// return reflect.Value{}, fmt.Errorf("parse int64 error: %v", err)
				}
				elem.Field(i).SetInt(f)
			case reflect.Uint64, reflect.Uint32:
				if s == "" {
					elem.Field(i).SetUint(0)
					break
				}

				f, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					ok = false
					break
					// return reflect.Value{}, fmt.Errorf("parse uint64 error: %v", err)
				}
				elem.Field(i).SetUint(f)
			case reflect.String:
				elem.Field(i).SetString(s)
			}
		}

		if ok {
			slice = append(slice, t)
		} else {
			dirty = append(dirty, data)
		}
	}

	return slice, dirty, nil
}
