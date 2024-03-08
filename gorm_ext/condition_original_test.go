package gorm_ext

import (
	"fmt"
	"testing"
)

func Test_OriginalCondition_error(t *testing.T) {
	//sql为空
	var cond = &OriginalCondition{
		Sql: "",
		Arg: "123",
	}
	_, _, err := cond.BuildSql(dbType)
	if err == nil {
		t.Errorf("Test_OriginalCondition_error faild")
	}
}

func Test_OriginalCondition_Eq(t *testing.T) {
	//1
	var cond = &OriginalCondition{
		Sql: "a = ?",
		Arg: "123",
	}
	sql, param, err := cond.BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_OriginalCondition_Eq faild")
	}
	if sql != fmt.Sprintf("a = ?") {
		t.Errorf(fmt.Sprintf("Test_OriginalCondition_Eq faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != "123" {
		t.Errorf("Test_OriginalCondition_Eq faild")
	}

}
