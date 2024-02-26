package gorm_ext

// Condition 定义查询条件
type Condition interface {
	GetSql() string           //获取当前条件的 sql
	GetParams() []interface{} //获取当前条件的 参数
}

type condition struct {
}

func (c *condition) GetSql() string {
	return ""
}

func (c *condition) GetParams() []interface{} {
	return nil
}
