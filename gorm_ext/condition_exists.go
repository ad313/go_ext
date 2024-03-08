package gorm_ext

import (
	"errors"
	"fmt"
	"gorm.io/gorm/schema"
)

// ExistsCondition Exists 和 Not Exists
type ExistsCondition struct {
	Table            schema.Tabler     //指定内表
	ConditionBuilder *ConditionBuilder //条件
	IsNotExists      bool              //默认false：exists；true ：not exists
	Func             string            //数据库函数，max、min 等，给当前字段套上函数

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

func (c *ExistsCondition) BuildSql(dbType string, extend ...interface{}) (string, []interface{}, error) {
	if !c.isBuild {
		if dbType == "" {
			c.error = errors.New("请指定数据库类型")
			return "", nil, c.error
		}

		if c.Table == nil {
			c.error = errors.New("请指定表")
			return "", nil, c.error
		}

		var isUnscoped = false
		if len(extend) > 0 {
			if v, ok := extend[0].(bool); ok {
				isUnscoped = v
			}
		}
		c.sql, c.params, c.error = c.buildExistsMethod(dbType, isUnscoped)
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

func (c *ExistsCondition) buildExistsMethod(dbType string, isUnscoped bool) (string, []interface{}, error) {
	if c == nil {
		return "", nil, errors.New("ExistsCondition 不能为空")
	}

	if c.ConditionBuilder == nil {
		return "", nil, errors.New("ExistsCondition Columns 不能为空")
	}

	var sql = fmt.Sprintf("SELECT 1 FROM %v WHERE ", formatSqlName(c.Table.TableName(), dbType))

	//处理软删除字段
	if !isUnscoped {
		softDel, err := getTableSoftDeleteColumnSql(c.Table, "", dbType)
		if err != nil {
			return "", nil, err
		}
		sql += softDel + " AND "
	}

	where, param, err := c.ConditionBuilder.BuildSql(dbType)
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
