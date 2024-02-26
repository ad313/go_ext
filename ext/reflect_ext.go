// Package ext 反射扩展
package ext

import (
	"reflect"
)

// IsType 判断 TSource 是否等于或者实现 Target
func IsType[TSource any, Target any]() (*TSource, *Target, bool) {
	var source = new(TSource)
	v, ok := IsTypeByValue[Target](*source)
	if ok {
		return source, v, true
	}

	return nil, nil, false
}

// IsTypeByValue 判断 value 是否等于或者实现 T
func IsTypeByValue[T any](value any) (*T, bool) {

	var inputRealValue interface{}

	//判断是否是指针
	if realValue, ok := IsPointerReturnValue(value); ok {
		inputRealValue = realValue
	} else {
		inputRealValue = value
	}

	if instance, ok := inputRealValue.(T); ok {
		return &instance, true
	}

	if instance, ok := reflect.New(reflect.TypeOf(inputRealValue)).Interface().(T); ok {
		return &instance, true
	}

	return nil, false
}

// IsInherit 判断 TSource 是否继承 Target
func IsInherit[TSource any, Target any]() bool {
	var source interface{} = *new(TSource)
	var target interface{} = *new(Target)

	v, ok := IsPointerReturnValue(source)
	if ok {
		source = v
	}

	v2, ok := IsPointerReturnValue(target)
	if ok {
		target = v2
	}

	targetType := reflect.TypeOf(target)
	sourceType := reflect.TypeOf(source)

	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		if field.PkgPath == targetType.PkgPath() && field.Name == targetType.Name() {
			return true
		}
	}

	return false
}

// IsPointer 判断是否是指针
func IsPointer(param interface{}) bool {
	return reflect.ValueOf(param).Kind() == reflect.Ptr
}

// IsPointerReturnValue 判断是否是指针，并返回真实值
func IsPointerReturnValue(param interface{}) (interface{}, bool) {
	value := reflect.ValueOf(param)
	if value.Kind() == reflect.Ptr {
		//处理 nil
		if value.Pointer() == 0 {
			//根据类型创建一个非nil的值，重新判断
			return IsPointerReturnValue(reflect.New(reflect.TypeOf(param).Elem()).Interface())
		}

		return value.Elem().Interface(), true
	}
	return nil, false
}

// GetPath 获取类型名称
func GetPath[T interface{}]() string {
	var value interface{}
	if IsPointer(*new(T)) {
		value = new(T)
	} else {
		value = (*T)(nil)
	}

	return reflect.TypeOf(value).Elem().String()
}
