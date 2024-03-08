package gorm_ext

import (
	"errors"
)

// Condition 表与值比较条件
type Condition struct {
	TableAlias     string      //表别名
	Column         any         //字段名
	CompareSymbols string      //比较符号
	Arg            interface{} //sql 参数
	Func           string      //数据库函数，max、min 等，给当前字段套上函数

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

func (c *Condition) BuildSql(dbType string, extend ...interface{}) (string, []interface{}, error) {
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

		c.sql, c.error = mergeWhereString(c.Column, c.CompareSymbols, c.TableAlias, c.Func, dbType)
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
