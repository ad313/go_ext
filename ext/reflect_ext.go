// Package ext 反射扩展
package ext

import (
	"reflect"
)

// IsType 判断两种类型是否 相等、继承
func IsType[TSource any, Target any]() (*TSource, *Target, bool) {
	var source = new(TSource)
	v, ok := IsTypeByValue[Target](*source)
	if ok {
		return source, v, true
	}

	return nil, nil, false
}

// IsTypeByValue IsType 判断给定的值和类型是否 相等、继承
func IsTypeByValue[T any](value any) (*T, bool) {

	var inputPointer = false
	var inputRealValue interface{}

	//判断是否是指针
	if realValue, ok := IsPointerReturnValue(value); ok {
		inputPointer = true
		inputRealValue = realValue
	} else {
		inputRealValue = value
	}

	//类型直接比较
	if IsPointer(*new(T)) {
		if inputPointer {
			if instance, ok := value.(T); ok {
				return &instance, true
			}
		} else {
			if instance, ok := value.(T); ok {
				return &instance, true
			}

			if instance, ok := reflect.New(reflect.TypeOf(value)).Interface().(T); ok {
				return &instance, true
			}
		}
	} else {
		if inputPointer {
			value = inputRealValue
			if instance, ok := value.(T); ok {
				return &instance, true
			}
		} else {
			if instance, ok := value.(T); ok {
				return &instance, true
			}
		}
	}

	//通过反射比较
	var formatValue interface{}
	if IsPointer(*new(T)) {
		if inputPointer {
			formatValue = value
		} else {
			formatValue = value
		}
	} else {
		if inputPointer {
			formatValue = inputRealValue
		} else {
			formatValue = value
		}
	}

	if instance, ok := reflect.New(reflect.TypeOf(formatValue)).Interface().(T); ok {
		return &instance, true
	}

	return nil, false
}

// IsPointer 判断是否是指针
func IsPointer(param interface{}) bool {
	return reflect.ValueOf(param).Kind() == reflect.Ptr
}

// IsPointerReturnValue 判断是否是指针，并返回真实值
func IsPointerReturnValue(param interface{}) (interface{}, bool) {
	value := reflect.ValueOf(param)
	if value.Kind() == reflect.Ptr {
		if value.Pointer() == 0 {
			return nil, false
		}

		return value.Elem().Interface(), true
	}
	return nil, false
}
