package gorm_ext

import (
	"fmt"
	"testing"
)

func Test_Condition_table_error(t *testing.T) {
	//字段为空 1
	var cond = &TableCondition{
		InnerAlias:     "",
		InnerColumn:    nil,
		OuterAlias:     "",
		OuterColumn:    nil,
		CompareSymbols: "",
	}
	_, _, err := cond.Build(dbType)
	if err == nil {
		t.Errorf("Test_Condition_table_error faild")
	}

	//字段为空 2
	cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "",
		OuterColumn:    nil,
		CompareSymbols: "",
	}
	_, _, err = cond.Build(dbType)
	if err == nil {
		t.Errorf("Test_Condition_table_error faild")
	}

	//字段为空 3
	cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    nil,
		CompareSymbols: "",
	}
	_, _, err = cond.Build(dbType)
	if err == nil {
		t.Errorf("Test_Condition_table_error faild")
	}

	//符号为空
	cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Name,
		CompareSymbols: "",
	}
	_, _, err = cond.Build(dbType)
	if err == nil {
		t.Errorf("Test_Condition_table_error faild")
	}

}

func Test_Condition_table_Eq(t *testing.T) {
	//1
	var cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Name,
		CompareSymbols: Eq,
	}
	sql, param, err := cond.clear().Build(dbType)
	if err != nil {
		t.Errorf("Test_Condition_table_Eq faild")
	}
	if sql != fmt.Sprintf("%v.%v = %v.%v", f("a"), f("id"), f("b"), f("name")) {
		t.Errorf(fmt.Sprintf("Test_Condition_table_Eq faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_Condition_table_Eq faild")
	}

	//2
	cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Name,
		CompareSymbols: NotEq,
	}
	sql, param, err = cond.clear().Build(dbType)
	if err != nil {
		t.Errorf("Test_Condition_table_Eq faild")
	}
	if sql != fmt.Sprintf("%v.%v <> %v.%v", f("a"), f("id"), f("b"), f("name")) {
		t.Errorf(fmt.Sprintf("Test_Condition_table_Eq faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_Condition_table_Eq faild")
	}

	//3
	cond = &TableCondition{
		InnerAlias:     "",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Name,
		CompareSymbols: NotEq,
	}
	sql, param, err = cond.clear().Build(dbType)
	if err != nil {
		t.Errorf("Test_Condition_table_Eq faild")
	}
	if sql != fmt.Sprintf("%v <> %v.%v", f("id"), f("b"), f("name")) {
		t.Errorf(fmt.Sprintf("Test_Condition_table_Eq faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_Condition_table_Eq faild")
	}
}

func Test_Condition_table_IsNull(t *testing.T) {
	//1
	var cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Name,
		CompareSymbols: IsNull,
	}
	sql, param, err := cond.clear().Build(dbType)
	if err != nil {
		t.Errorf("Test_Condition_table_IsNull faild")
	}
	if sql != fmt.Sprintf("%v.%v IS NULL ", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_table_IsNull faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_Condition_table_IsNull faild")
	}

	//2
	cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Name,
		CompareSymbols: NotNull,
	}
	sql, param, err = cond.clear().Build(dbType)
	if err != nil {
		t.Errorf("Test_Condition_table_IsNull faild")
	}
	if sql != fmt.Sprintf("%v.%v IS NOT NULL ", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_table_IsNull faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_Condition_table_IsNull faild")
	}
}
