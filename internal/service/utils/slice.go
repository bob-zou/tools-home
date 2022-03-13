package utils

import "reflect"

// SliceContains 判断元素是否在切片中
func SliceContains(sli interface{}, elem interface{}) bool {
	sliValue := reflect.ValueOf(sli)
	if sliValue.Kind() != reflect.Slice && sliValue.Kind() != reflect.Array {
		return false
	}

	if sliValue.Len() == 0 {
		return false
	}

	elemValue := reflect.ValueOf(elem)
	if elemValue.Kind() == reflect.Struct {
		// 不支持结构体
		return false
	}

	if sliValue.Index(0).Kind() != elemValue.Kind() {
		return false
	}

	for i := 0; i < sliValue.Len(); i++ {
		if sliValue.Index(i).Interface() == elem {
			return true
		}
	}

	return false
}
