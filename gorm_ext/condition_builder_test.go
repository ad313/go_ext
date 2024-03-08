package gorm_ext

import (
	"fmt"
	"testing"
)

func Test_gorm_condition_builder_error(t *testing.T) {
	//字段为空 1
	var cond = &TableCondition{
		InnerAlias:     "",
		InnerColumn:    nil,
		OuterAlias:     "",
		OuterColumn:    nil,
		CompareSymbols: "",
	}
	_, _, err := NewAndEmptyConditionBuilder().AddChildrenCondition(cond).BuildSql(dbType)
	if err == nil {
		t.Errorf("Test_gorm_condition_builder faild")
	}
}

func Test_gorm_condition_builder(t *testing.T) {
	//1
	var cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Name,
		CompareSymbols: NotEq,
	}
	sql, param, err := NewAndEmptyConditionBuilder().AddChildrenCondition(cond).BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_gorm_condition_builder faild")
	}
	if sql != fmt.Sprintf("(%v.%v <> %v.%v)", f("a"), f("id"), f("b"), f("name")) {
		t.Errorf(fmt.Sprintf("Test_gorm_condition_builder faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_gorm_condition_builder faild")
	}

	//2 and
	var cond2 = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Age,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Age,
		CompareSymbols: Gt,
	}
	sql, param, err = NewAndEmptyConditionBuilder().AddChildrenCondition(cond, cond2).BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_gorm_condition_builder faild")
	}
	if sql != fmt.Sprintf("(%v.%v <> %v.%v AND %v.%v > %v.%v)", f("a"), f("id"), f("b"), f("name"), f("a"), f("age"), f("b"), f("age")) {
		t.Errorf(fmt.Sprintf("Test_gorm_condition_builder faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_gorm_condition_builder faild")
	}

	//3 or
	sql, param, err = NewOrEmptyConditionBuilder().AddChildrenCondition(cond, cond2).BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_gorm_condition_builder faild")
	}
	if sql != fmt.Sprintf("(%v.%v <> %v.%v OR %v.%v > %v.%v)", f("a"), f("id"), f("b"), f("name"), f("a"), f("age"), f("b"), f("age")) {
		t.Errorf(fmt.Sprintf("Test_gorm_condition_builder faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_gorm_condition_builder faild")
	}

	//4 and
	var cond_IsNull = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Age,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Age,
		CompareSymbols: IsNull, //此时 OuterAlias 和 OuterColumn 无效
	}
	sql, param, err = NewAndEmptyConditionBuilder().AddChildrenCondition(cond_IsNull).BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_gorm_condition_builder faild")
	}
	if sql != fmt.Sprintf("(%v.%v IS NULL )", f("a"), f("age")) {
		t.Errorf(fmt.Sprintf("Test_gorm_condition_builder faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_gorm_condition_builder faild")
	}

	//嵌套
	var b1 = NewOrConditionBuilder(cond, cond2)
	sql, param, err = NewAndEmptyConditionBuilder().AddChildrenBuilder(b1).AddChildrenCondition(cond_IsNull).BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_gorm_condition_builder faild")
	}

	if sql != "((\"a\".\"id\" <> \"b\".\"name\" OR \"a\".\"age\" > \"b\".\"age\") AND \"a\".\"age\" IS NULL )" {
		t.Errorf(fmt.Sprintf("Test_gorm_condition_builder faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_gorm_condition_builder faild")
	}
}
