package utils

import (
	"math"
	"reflect"
	"unsafe"
)

// TrunFloat 截断小数
func TrunFloat(f float64, prec int) float64 {
	x := math.Pow10(prec)
	return math.Trunc(f*x) / x
}

// DivFloat 除法
func DivFloat(f1, f2 float64, prec int) float64 {
	if f2 == 0 {
		return 0
	}
	return TrunFloat(f1/f2, prec)
}

// S2b 字符串转字节数组
func S2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// B2s 字节数组转字符串
func B2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func IsNil(a any) bool {
	if a == nil {
		return true
	}
	return reflect.ValueOf(a).IsNil()
}

// Ptr 返回一个指针
func Ptr[T any](t T) *T {
	return &t
}

// Ptr2 返回一个指针，用于简化代码, 根据returnSelf判断是否返回空指针
func Ptr2[T any](t T, returnSelf bool) *T {
	if returnSelf {
		return &t
	}
	return nil
}

func Value[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}
