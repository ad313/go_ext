package gorm_ext

import (
	"gorm.io/gorm"
	"testing"
)

type Table1 struct {
	Id     string `gorm:"column:id;type:VARCHAR2(36);primaryKey;not null"` //标识
	Table2        //基类
}

type Table2 struct {
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME(8)"` //删除标识
}

var tableName = "table1"

func (t *Table1) TableName() string {
	return tableName
}

func Test_BuildGormTable(t *testing.T) {
	var result = BuildGormTable[Table1]()
	if result == nil || result.Table == nil || result.Error != nil {
		t.Errorf("BuildGormTable faild")
		return
	}

	if result.Table.Name != tableName {
		t.Errorf("BuildGormTable faild")
		return
	}

	if result.Table.DeletedColumnName != "deleted_at" {
		t.Errorf("BuildGormTable faild")
		return
	}

	//缓存
	result = BuildGormTable[Table1]()
	if result == nil || result.Table == nil || result.Error != nil {
		t.Errorf("BuildGormTable faild")
		return
	}

	if result.Table.Name != tableName {
		t.Errorf("BuildGormTable faild")
		return
	}

	if result.Table.DeletedColumnName != "deleted_at" {
		t.Errorf("BuildGormTable faild")
		return
	}
}
