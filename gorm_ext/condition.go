package gorm_ext

// Condition 常规条件
type Condition struct {
	TableAlias     string      //表别名
	Column         any         //字段名
	CompareSymbols string      //比较符号
	Arg            interface{} //参数
	error          error       //错误
}

func (c *Condition) GetSql() string {
	return ""
}

func (c *Condition) GetParams() []interface{} {
	if c.Arg == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{c.Arg}
}

func (c *Condition) GetError() error {
	return c.error
}
