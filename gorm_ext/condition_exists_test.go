package gorm_ext

import (
	"fmt"
	"testing"
)

func Test_condition_exists_error(t *testing.T) {
	//字段为空 1
	var cond = &ExistsCondition{
		Table:       nil,
		Column:      nil,
		IsNotExists: false,
	}
	_, _, err := cond.Build(dbType)
	if err == nil {
		t.Errorf("Test_condition_exists_error faild")
	}

	//2
	cond = &ExistsCondition{
		Table:       condTable,
		Column:      nil,
		IsNotExists: false,
	}
	_, _, err = cond.Build(dbType)
	if err == nil {
		t.Errorf("Test_condition_exists_error faild")
	}
}

func Test_condition_exists(t *testing.T) {
	//1
	var cond = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Id,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Name,
		CompareSymbols: NotEq,
	}

	//2 and
	var cond2 = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Age,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Age,
		CompareSymbols: Gt,
	}

	//4 and
	var cond_IsNull = &TableCondition{
		InnerAlias:     "a",
		InnerColumn:    &condTable.Age,
		OuterAlias:     "b",
		OuterColumn:    &condTable.Age,
		CompareSymbols: IsNull, //此时 OuterAlias 和 OuterColumn 无效
	}

	var exists = ExistsCondition{
		Table:       condTable,
		Column:      nil,
		IsNotExists: false,
	}

	exists.Column = NewOrEmptyConditionBuilder().AddChildrenBuilder(NewAndConditionBuilder(cond, cond2)).AddChildrenCondition(cond_IsNull)
	sql, param, err := exists.Build(dbType)
	if err != nil {
		t.Errorf("Test_condition_exists faild")
	}

	var targetSql = ""
	switch dbType {
	case Dm:
		targetSql = "Exists (SELECT 1 FROM \"conditionTable\" WHERE ((\"a\".\"id\" <> \"b\".\"name\" AND \"a\".\"age\" > \"b\".\"age\") OR \"a\".\"age\" IS NULL ))"
	}

	if sql != targetSql {
		t.Errorf(fmt.Sprintf("Test_condition_exists faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_condition_exists faild")
	}

	//Not exists
	exists.IsNotExists = true
	sql, param, err = exists.clear().Build(dbType)
	if err != nil {
		t.Errorf("Test_condition_exists faild")
	}

	switch dbType {
	case Dm:
		targetSql = "Not Exists (SELECT 1 FROM \"conditionTable\" WHERE ((\"a\".\"id\" <> \"b\".\"name\" AND \"a\".\"age\" > \"b\".\"age\") OR \"a\".\"age\" IS NULL ))"
	}

	if sql != targetSql {
		t.Errorf(fmt.Sprintf("Test_condition_exists faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_condition_exists faild")
	}

}
