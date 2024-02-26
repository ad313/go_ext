package ext

import (
	"gorm.io/gorm/schema"
	"testing"
)

type table struct {
	Id int32
	table2
}

type table2 struct {
	Name string
}

func (a *table) TableName() string {
	return ""
}

var model = table{Id: 1}
var modelPointer = &model

// Test_IsType 判断 TSource 是否等于或者实现 Target
func Test_IsType(t *testing.T) {
	var ok = false

	//实现接口
	_, _, ok = IsType[table, schema.Tabler]()
	if !ok {
		t.Errorf("IsType faild")
	}

	_, _, ok = IsType[*table, schema.Tabler]()
	if !ok {
		t.Errorf("IsType faild")
	}

	//类型本身
	_, _, ok = IsType[table, table]()
	if !ok {
		t.Errorf("IsType faild")
	}

	_, _, ok = IsType[*table, table]()
	if !ok {
		t.Errorf("IsType faild")
	}

	_, _, ok = IsType[table, *table]()
	if !ok {
		t.Errorf("IsType faild")
	}

	_, _, ok = IsType[*table, *table]()
	if !ok {
		t.Errorf("IsType faild")
	}
}

// Test_IsTypeByValue 判断 value 是否等于或者实现 T
func Test_IsTypeByValue(t *testing.T) {
	var ok = false

	//类型本身
	_, ok = IsTypeByValue[table](model)
	if !ok {
		t.Errorf("IsTypeByValue faild")
	}
	_, ok = IsTypeByValue[table](modelPointer)
	if !ok {
		t.Errorf("IsTypeByValue faild")
	}
	_, ok = IsTypeByValue[*table](model)
	if !ok {
		t.Errorf("IsTypeByValue faild")
	}
	_, ok = IsTypeByValue[*table](modelPointer)
	if !ok {
		t.Errorf("IsTypeByValue faild")
	}

	//实现接口
	_, ok = IsTypeByValue[schema.Tabler](model)
	if !ok {
		t.Errorf("IsTypeByValue faild")
	}
	_, ok = IsTypeByValue[schema.Tabler](modelPointer)
	if !ok {
		t.Errorf("IsTypeByValue faild")
	}
}

// Test_IsPointer 是否是指针
func Test_IsPointer(t *testing.T) {
	ok := IsPointer(model)
	if ok {
		t.Errorf("IsPointer faild")
	}

	var table *table = nil
	ok = IsPointer(table)
	if !ok {
		t.Errorf("IsPointer faild")
	}

	ok = IsPointer(modelPointer)
	if !ok {
		t.Errorf("IsPointer faild")
	}
}

// Test_IsPointerReturnValue 是否是指针，并返回真实值
func Test_IsPointerReturnValue(t *testing.T) {

	var target = table{Id: 1}

	_, ok := IsPointerReturnValue(target)
	if ok {
		t.Errorf("IsPointerReturnValue faild")
	}

	v, ok := IsPointerReturnValue(&target)
	if !ok {
		t.Errorf("IsPointerReturnValue faild")
	}

	if v == nil {
		t.Errorf("IsPointerReturnValue faild")
	}

	realValue, ok := IsTypeByValue[table](v)
	if !ok {
		t.Errorf("IsPointerReturnValue faild")
	}

	if realValue == nil || realValue.Id != target.Id {
		t.Errorf("IsPointerReturnValue faild")
	}
}

// Test_Inherit 是否继承
func Test_Inherit(t *testing.T) {
	var ok = false

	//类型本身
	ok = IsInherit[table, table2]()
	if !ok {
		t.Errorf("IsInherit faild")
	}

	ok = IsInherit[*table, table2]()
	if !ok {
		t.Errorf("IsType faild")
	}

	ok = IsInherit[table, *table2]()
	if !ok {
		t.Errorf("IsType faild")
	}

	ok = IsInherit[*table, *table2]()
	if !ok {
		t.Errorf("IsType faild")
	}
}

func Test_NotSupport(t *testing.T) {

	var ok = false

	//实现接口不支持指针
	_, ok = IsTypeByValue[*schema.Tabler](model)
	if ok {
		t.Errorf("IsTypeByValue faild")
	}
	_, ok = IsTypeByValue[*schema.Tabler](modelPointer)
	if ok {
		t.Errorf("IsTypeByValue faild")
	}

	//实现接口不支持指针
	_, _, ok = IsType[table, *schema.Tabler]()
	if ok {
		t.Errorf("IsType faild")
	}

	_, _, ok = IsType[*table, *schema.Tabler]()
	if ok {
		t.Errorf("IsType faild")
	}
}

func Test_GetPath(t *testing.T) {
	value := GetPath[table]()
	if value != "ext.table" {
		t.Errorf("GetPath faild")
	}

	value = GetPath[table2]()
	if value != "ext.table2" {
		t.Errorf("GetPath faild")
	}

	value = GetPath[*table]()
	if value != "ext.table" {
		t.Errorf("GetPath faild")
	}

	value = GetPath[*table2]()
	if value != "ext.table2" {
		t.Errorf("GetPath faild")
	}
}
