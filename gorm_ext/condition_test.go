package gorm_ext

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
)

type conditionTable struct {
	Id        string         `gorm:"column:id;type:VARCHAR2(36);primaryKey;not null"`
	Name      string         `gorm:"column:name;type:VARCHAR2(36);not null"`
	Age       int            `gorm:"column:age;type:INTEGER(4);not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME(8)" json:"deleted_at"` //删除标识
}

func (t *conditionTable) TableName() string {
	return "conditionTable"
}

var condTable = BuildGormTable[conditionTable]().Table.T

var dbType = Dm

func Test_Condition_error(t *testing.T) {
	//字段为空
	var cond = &Condition{
		TableAlias:     "a",
		Column:         "",
		CompareSymbols: Eq,
		Arg:            "123",
	}
	_, _, err := cond.BuildSql(dbType)
	if err == nil {
		t.Errorf("Test_Condition_error faild")
	}

	//符号为空
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: "",
		Arg:            "123",
	}
	_, _, err = cond.BuildSql(dbType)
	if err == nil {
		t.Errorf("Test_Condition_error faild")
	}

	//参数为空 1
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: Eq,
		Arg:            nil,
	}
	_, _, err = cond.BuildSql(dbType)
	if err == nil {
		t.Errorf("Test_Condition_error faild")
	}

	//参数为空 2
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: IsNull,
		Arg:            nil,
	}
	_, _, err = cond.BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_error faild")
	}

	//参数为空 3
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: NotNull,
		Arg:            nil,
	}
	_, _, err = cond.BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_error faild")
	}
}

func Test_Condition_Eq(t *testing.T) {
	//1
	var cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: Eq,
		Arg:            "123",
	}
	sql, param, err := cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Eq faild")
	}
	if sql != fmt.Sprintf("%v.%v = ?", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Eq faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != "123" {
		t.Errorf("Test_Condition_Eq faild")
	}

	//2
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: NotEq,
		Arg:            "123",
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Eq faild")
	}
	if sql != fmt.Sprintf("%v.%v <> ?", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Eq faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != "123" {
		t.Errorf("Test_Condition_Eq faild")
	}

	//3 func
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: NotEq,
		Arg:            "123",
		Func:           Count,
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Eq faild")
	}
	if sql != fmt.Sprintf("Count(%v.%v) <> ?", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Eq faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != "123" {
		t.Errorf("Test_Condition_Eq faild")
	}
}

func Test_Condition_Like(t *testing.T) {
	//1
	var cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: Like,
		Arg:            "123",
	}
	sql, param, err := cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Like faild")
	}
	if sql != fmt.Sprintf("%v.%v LIKE ?", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Like faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != "%123%" {
		t.Errorf("Test_Condition_Like faild")
	}

	//2
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: NotLike,
		Arg:            "123",
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Like faild")
	}
	if sql != fmt.Sprintf("%v.%v NOT LIKE ?", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Like faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != "%123%" {
		t.Errorf("Test_Condition_Like faild")
	}

	//3
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: StartWith,
		Arg:            "123",
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Like faild")
	}
	if sql != fmt.Sprintf("%v.%v LIKE ?", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Like faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != "123%" {
		t.Errorf("Test_Condition_Like faild")
	}

	//4
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: EndWith,
		Arg:            "123",
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Like faild")
	}
	if sql != fmt.Sprintf("%v.%v LIKE ?", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Like faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != "%123" {
		t.Errorf("Test_Condition_Like faild")
	}
}

func Test_Condition_In(t *testing.T) {
	//1
	var cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: In,
		Arg:            []string{"123", "456"},
	}
	sql, param, err := cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_In faild")
	}
	if sql != fmt.Sprintf("%v.%v IN (?)", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_In faild：a：%v", sql))
	}
	if len(param) != 1 {
		t.Errorf("Test_Condition_In faild")
	}

	//2
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: NotIn,
		Arg:            []string{"123", "456"},
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_In faild")
	}
	if sql != fmt.Sprintf("%v.%v NOT IN (?)", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_In faild：a：%v", sql))
	}
	if len(param) != 1 {
		t.Errorf("Test_Condition_In faild")
	}
}

func Test_Condition_IsNull(t *testing.T) {
	//1
	var cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: IsNull,
		Arg:            []string{"123", "456"},
	}
	sql, param, err := cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_IsNull faild")
	}
	if sql != fmt.Sprintf("%v.%v IS NULL ", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_IsNull faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_Condition_IsNull faild")
	}

	//2
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Id,
		CompareSymbols: NotNull,
		Arg:            []string{"123", "456"},
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_IsNull faild")
	}
	if sql != fmt.Sprintf("%v.%v IS NOT NULL ", f("a"), f("id")) {
		t.Errorf(fmt.Sprintf("Test_Condition_IsNull faild：a：%v", sql))
	}
	if len(param) != 0 {
		t.Errorf("Test_Condition_IsNull faild")
	}
}

func Test_Condition_Gt(t *testing.T) {
	//1
	var cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Age,
		CompareSymbols: Gt,
		Arg:            18,
	}
	sql, param, err := cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Gt faild")
	}
	if sql != fmt.Sprintf("%v.%v > ?", f("a"), f("age")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Gt faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != 18 {
		t.Errorf("Test_Condition_Gt faild")
	}

	//2
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Age,
		CompareSymbols: GtAndEq,
		Arg:            18,
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Gt faild")
	}
	if sql != fmt.Sprintf("%v.%v >= ?", f("a"), f("age")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Gt faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != 18 {
		t.Errorf("Test_Condition_Gt faild")
	}
}

func Test_Condition_Less(t *testing.T) {
	//1
	var cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Age,
		CompareSymbols: Less,
		Arg:            18,
	}
	sql, param, err := cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Less faild")
	}
	if sql != fmt.Sprintf("%v.%v < ?", f("a"), f("age")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Less faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != 18 {
		t.Errorf("Test_Condition_Less faild")
	}

	//2
	cond = &Condition{
		TableAlias:     "a",
		Column:         &condTable.Age,
		CompareSymbols: LessAndEq,
		Arg:            18,
	}
	sql, param, err = cond.clear().BuildSql(dbType)
	if err != nil {
		t.Errorf("Test_Condition_Less faild")
	}
	if sql != fmt.Sprintf("%v.%v <= ?", f("a"), f("age")) {
		t.Errorf(fmt.Sprintf("Test_Condition_Less faild：a：%v", sql))
	}
	if len(param) != 1 || param[0] != 18 {
		t.Errorf("Test_Condition_Less faild")
	}
}

// 处理数据库字段名称
func f(str string) string {
	return getSqlSm(dbType) + str + getSqlSm(dbType)
}
