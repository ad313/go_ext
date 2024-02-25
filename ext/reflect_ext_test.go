package ext

import (
	"gorm.io/gorm/schema"
	"testing"
)

type table struct {
	Id int32
}

func (a *table) TableName() string {
	return ""
}

var model = table{Id: 1}
var modelPointer = &model

func Test_IsType(t *testing.T) {
	var ok = false
	_, _, ok = IsType[table, schema.Tabler]()
	if !ok {
		t.Errorf("IsType faild")
	}

	_, _, ok = IsType[*table, schema.Tabler]()
	if !ok {
		t.Errorf("IsType faild")
	}

	//继承接口不支持指针
	_, _, ok = IsType[table, *schema.Tabler]()
	if ok {
		t.Errorf("IsType faild")
	}

	_, _, ok = IsType[*table, *schema.Tabler]()
	if ok {
		t.Errorf("IsType faild")
	}
}

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

	//继承接口
	_, ok = IsTypeByValue[schema.Tabler](model)
	if !ok {
		t.Errorf("IsTypeByValue faild")
	}
	_, ok = IsTypeByValue[schema.Tabler](modelPointer)
	if !ok {
		t.Errorf("IsTypeByValue faild")
	}

	//继承接口不支持指针
	_, ok = IsTypeByValue[*schema.Tabler](model)
	if ok {
		t.Errorf("IsTypeByValue faild")
	}
	_, ok = IsTypeByValue[*schema.Tabler](modelPointer)
	if ok {
		t.Errorf("IsTypeByValue faild")
	}
}

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
