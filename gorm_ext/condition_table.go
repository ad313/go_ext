package gorm_ext

import (
	"errors"
	"github.com/ad313/go_ext/ext"
)

// TableCondition 表与表之间比较条件
type TableCondition struct {
	InnerAlias     string //内部表表别名(exists)；左边表别名
	InnerColumn    any    //内部表字段名(exists)；左边表字段名
	OuterAlias     string //外部表表别名(exists)；右边表别名
	OuterColumn    any    //外部表字段名(exists)；右边表字段名
	CompareSymbols string //比较运算符号

	isBuild bool   //是否已经build
	sql     string //生成的sql
	//params  []interface{} //sql 参数
	error error //错误
}

func (c *TableCondition) Build(dbType string) (string, []interface{}, error) {
	if !c.isBuild {
		if dbType == "" {
			c.error = errors.New("请指定数据库类型")
			return "", nil, c.error
		}

		if c.CompareSymbols == "" {
			c.error = errors.New("比较运算符号不能为空")
			return "", nil, c.error
		}

		//左边sql
		innerSql, err := mergeTableColumn(c.InnerColumn, c.InnerAlias, dbType)
		c.error = err
		if c.error != nil {
			return "", nil, c.error
		}

		//判断是否需要右边的参数，固定传1，再比较
		flag, err := checkParam(c.CompareSymbols, "1")
		c.error = err
		if c.error != nil {
			return "", nil, c.error
		}

		if flag == "1" {
			//右边sql
			outerSql, err := mergeTableColumn(c.OuterColumn, c.OuterAlias, dbType)
			c.error = err
			if c.error != nil {
				return "", nil, c.error
			}

			c.sql = innerSql + " " + getCompareSymbols(c.CompareSymbols) + " " + outerSql
		} else {
			c.sql = innerSql + " " + getCompareSymbols(c.CompareSymbols) + " "
		}

		c.isBuild = true
		return c.sql, nil, nil
	}
	return c.sql, nil, c.error
}

func (c *TableCondition) clear() *TableCondition {
	if c.isBuild {
		c.isBuild = false
		c.sql = ""
		c.error = nil
	}

	return c
}

func mergeTableColumn(column interface{}, table string, dbType string) (string, error) {
	name, err := resolveColumnName(column, dbType)
	if err != nil {
		return "", err
	}

	return ext.ChooseTrueValue(table != "", formatTableAlias(table, dbType)+"."+name, name), nil
}
