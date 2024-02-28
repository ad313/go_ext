package gorm_ext

import "gorm.io/gorm/schema"

// ExistsCondition Exists 和 Not Exists
type ExistsCondition struct {
	Table       *schema.Tabler
	Columns     []*ExistsColumn
	IsNotExists bool
	error       error //错误
}

type ExistsConditionDetail struct {
	InnerColumn    any         //exists 对应的表字段
	OuterColumn    any         //外部表字段
	OuterAlias     string      //外部表表别名
	OuterValue     interface{} //外部直接传值
	CompareSymbols string      //比较符号
}

func (c *ExistsCondition) GetSql() string {
	return ""
}

func (c *ExistsCondition) GetParams() []interface{} {
	return nil
}

func (c *ExistsCondition) GetError() error {
	return c.error
}
