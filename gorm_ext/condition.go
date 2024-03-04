package gorm_ext

import (
	"errors"
	"fmt"
	"reflect"
)

// Condition 表与值比较条件
type Condition struct {
	TableAlias     string      //表别名
	Column         any         //字段名
	CompareSymbols string      //比较符号
	Arg            interface{} //sql 参数

	isBuild bool   //是否已经build
	sql     string //生成的sql
	//params  []interface{} //sql 参数
	error error //错误
}

func (c *Condition) getParams() []interface{} {

	if c.Arg == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{c.Arg}
}

func (c *Condition) Build(dbType string) (string, []interface{}, error) {
	if !c.isBuild {
		if dbType == "" {
			c.error = errors.New("请指定数据库类型")
			return "", nil, c.error
		}

		//检查参数有效性
		param, err := checkParam(c.CompareSymbols, c.Arg)
		if err != nil {
			c.error = err
			return "", nil, c.error
		}
		c.Arg = param

		c.sql, c.error = mergeWhereString(c.Column, c.CompareSymbols, c.TableAlias, dbType)
		c.Arg = mergeWhereValue(c.CompareSymbols, c.Arg)
		c.isBuild = true
	}
	return c.sql, c.getParams(), c.error
}

func (c *Condition) clear() *Condition {
	if c.isBuild {
		c.isBuild = false
		c.sql = ""
		c.error = nil
	}

	return c
}

// mergeWhereString 组合 where 条件
func mergeWhereString(column any, compareSymbols string, tableAlias string, dbType string) (string, error) {
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
		table = formatTableAlias(tableAlias, dbType) + "."
	}

	return fmt.Sprintf("%v%v %v %v", table, name, getCompareSymbols(compareSymbols), valueExpress), nil
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

func formatTableAlias(alias string, dbType string) string {
	if alias == "" {
		return alias
	}

	return getSqlSm(dbType) + alias + getSqlSm(dbType)
}
