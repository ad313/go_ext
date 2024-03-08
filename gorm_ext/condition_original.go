package gorm_ext

import "errors"

// OriginalCondition gorm 原始where条件,不作任何处理
type OriginalCondition struct {
	Sql string      //原始sql，不经过任何处理
	Arg interface{} //sql 参数

	isBuild bool  //是否已经build
	error   error //错误
}

func (c *OriginalCondition) getParams() []interface{} {

	if c.Arg == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{c.Arg}
}

func (c *OriginalCondition) BuildSql(dbType string, extend ...interface{}) (string, []interface{}, error) {
	if !c.isBuild {
		//if dbType == "" {
		//	c.error = errors.New("请指定数据库类型")
		//	return "", nil, c.error
		//}

		////检查参数有效性
		//param, err := checkParam(c.CompareSymbols, c.Arg)
		//if err != nil {
		//	c.error = err
		//	return "", nil, c.error
		//}
		//c.Arg = param

		//c.sql, c.error = mergeWhereString(c.Column, c.CompareSymbols, c.TableAlias, dbType)
		//c.Arg = mergeWhereValue(c.CompareSymbols, c.Arg)

		if c.Sql == "" {
			c.error = errors.New("sql不能为空")
			return "", nil, c.error
		}

		c.isBuild = true
	}
	return c.Sql, c.getParams(), c.error
}

func (c *OriginalCondition) clear() *OriginalCondition {
	if c.isBuild {
		c.isBuild = false
		c.error = nil
		c.Sql = ""
		c.Arg = nil
	}

	return c
}
