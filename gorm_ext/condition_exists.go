package gorm_ext

import (
	"errors"
	"fmt"
	"gorm.io/gorm/schema"
)

// ExistsCondition Exists 和 Not Exists
type ExistsCondition struct {
	Table       schema.Tabler
	Column      *ConditionBuilder
	IsNotExists bool

	isBuild bool          //是否已经build
	sql     string        //生成的sql
	params  []interface{} //sql 参数
	error   error         //错误
}

func (c *ExistsCondition) getParams() []interface{} {
	if c.params == nil {
		return make([]interface{}, 0)
	}
	return c.params
}

func (c *ExistsCondition) Build(dbType string) (string, []interface{}, error) {
	if !c.isBuild {
		if dbType == "" {
			c.error = errors.New("请指定数据库类型")
			return "", nil, c.error
		}

		if c.Table == nil {
			c.error = errors.New("请指定表")
			return "", nil, c.error
		}

		c.sql, c.params, c.error = c.buildExistsMethod(dbType)
		c.isBuild = true
	}

	return c.sql, c.getParams(), c.error
}

func (c *ExistsCondition) clear() *ExistsCondition {
	if c.isBuild {
		c.isBuild = false
		c.sql = ""
		c.error = nil
		c.params = nil
	}

	return c
}

func (c *ExistsCondition) buildExistsMethod(dbType string) (string, []interface{}, error) {
	if c == nil {
		return "", nil, errors.New("ExistsCondition 不能为空")
	}

	if c.Column == nil {
		return "", nil, errors.New("ExistsCondition Columns 不能为空")
	}

	var sql = fmt.Sprintf("SELECT 1 FROM %v WHERE ", formatTableAlias(c.Table.TableName(), dbType))

	where, param, err := c.Column.BuildSql(dbType)
	if err != nil {
		return "", nil, err
	}
	sql += where

	var first = "Exists"
	if c.IsNotExists {
		first = "Not " + first
	}
	return fmt.Sprintf("%v (%v)", first, sql), param, nil
}
