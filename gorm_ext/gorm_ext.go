package gorm_ext

import (
	"errors"
	"fmt"
	"github.com/ad313/go_ext/ext"
	"gorm.io/gorm/schema"
	"reflect"
)

// ResolveTableColumnName 从缓存获取数据库字段名称：如果不能匹配，则返回 string 值
func ResolveTableColumnName(column any, dbType string) string {
	var kind = reflect.ValueOf(column).Kind()
	if kind == reflect.Pointer {
		var name = GetTableColumn(column)
		if name == "" {
			return ""
		}
		return getSqlSm(dbType) + name + getSqlSm(dbType)
	} else {
		if str, ok := column.(string); ok && str != "" {
			return getSqlSm(dbType) + str + getSqlSm(dbType)
		} else {
			return ""
		}
	}
}

// mergeWhereString 组合 where 条件
func mergeWhereString(column any, compareSymbols string, tableAlias string, f string, dbType string) (string, error) {
	name, err := resolveColumnName(column, dbType)
	if err != nil {
		return "", err
	}

	var valueExpress = "?"
	switch compareSymbols {
	case "IN":
		valueExpress = "(?)"
		break
	case "NOT IN":
		valueExpress = "(?)"
		break
	case "IS NULL":
		valueExpress = ""
		break
	case "IS NOT NULL":
		valueExpress = ""
		break
	}

	var table = ""
	if tableAlias != "" {
		table = formatSqlName(tableAlias, dbType) + "."
	}

	name = table + name
	return fmt.Sprintf("%v %v %v", mergeNameAndFunc(name, f), getCompareSymbols(compareSymbols), valueExpress), nil
}

// mergeWhereValue 处理查询值，当 IN 条件时，拼接 %
func mergeWhereValue(compareSymbols string, value interface{}) interface{} {
	if value == nil {
		return value
	}
	v, ok := value.(string)
	if !ok {
		return value
	}

	switch compareSymbols {
	case "LIKE":
		return "%" + v + "%"
	case "NOT LIKE":
		return "%" + v + "%"
	case "STARTWITH":
		return v + "%"
	case "ENDWITH":
		return "%" + v
	}

	return value
}

// 检查各种条件下参数是否为空
func checkParam(compareSymbols string, value interface{}) (interface{}, error) {
	if compareSymbols == "" {
		return nil, errors.New("compareSymbols 不能为空")
	}

	var check = true
	switch compareSymbols {
	case "LIKE":
		break
	case "NOT LIKE":
		break
	case "STARTWITH":
		break
	case "ENDWITH":
		break
	case "IN":
		break
	case "NOT IN":
		break
	case "IS NULL":
		check = false
		break
	case "IS NOT NULL":
		check = false
		break
	}

	//不需要参数
	if !check {
		value = nil
	}

	if check && value == nil {
		return nil, errors.New("参数不能为空")
	}

	return value, nil
}

func getCompareSymbols(compareSymbols string) string {
	switch compareSymbols {
	case "LIKE":
		return compareSymbols
	case "NOT LIKE":
		return compareSymbols
	case "STARTWITH":
		return "LIKE"
	case "ENDWITH":
		return "LIKE"
	}

	return compareSymbols
}

// resolveColumnName 从缓存获取数据库字段名称：如果不能匹配，则返回 string 值
func resolveColumnName(column any, dbType string) (string, error) {
	var name = ResolveTableColumnName(column, dbType)
	if name == "" {
		return "", errors.New("未获取到字段名称")
	}
	return name, nil
}

// 处理数据库表名
func formatSqlName(alias string, dbType string) string {
	if alias == "" {
		return alias
	}

	return getSqlSm(dbType) + alias + getSqlSm(dbType)
}

// 处理数据库表名 加上别名
func mergeTableWithAlias(table string, alias string, dbType string) string {
	if table == "" {
		return table
	}

	table = getSqlSm(dbType) + table + getSqlSm(dbType)

	if alias != "" {
		table += " as " + alias
	}

	return table
}

// 处理数据库表名 加上别名
func mergeTableWithAliasByValue(table schema.Tabler, alias string, dbType string) string {
	return mergeTableWithAlias(table.TableName(), alias, dbType)
}

// 获取软删除字段
func getTableSoftDeleteColumnSql(table schema.Tabler, tableAlias string, dbType string) (string, error) {
	var tableSchema = GetTableSchema(table)
	if tableSchema != nil && tableSchema.DeletedColumnName != "" {
		n, err := resolveColumnName(tableSchema.DeletedColumnName, dbType)
		if err != nil {
			return "", err
		}

		if tableAlias != "" {
			n = formatSqlName(tableAlias, _dbType) + "." + n
		}

		return n + " IS NULL", nil
	}

	return "", nil
}

func mergeTableColumnWithFunc(column interface{}, table string, f string, dbType string) (string, error) {
	name, err := resolveColumnName(column, dbType)
	if err != nil {
		return "", err
	}

	return ext.ChooseTrueValue(table != "", mergeNameAndFunc(formatSqlName(table, dbType)+"."+name, f), mergeNameAndFunc(name, f)), nil
}

// 合并字段和数据库函数
func mergeNameAndFunc(name, f string) string {
	return ext.ChooseTrueValue(f == "", name, f+"("+name+")")
}
