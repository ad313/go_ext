package gorm_ext

import "errors"

// WhereCondition 定义查询条件
type WhereCondition interface {
	GetSql() string           //获取当前条件的 sql
	GetParams() []interface{} //获取当前条件的 参数
	GetError() error          //获取生成sql时的错误
}

// ConditionBuilder 条件构建器
type ConditionBuilder struct {
	Or      bool                //and、or
	Items   []*ConditionBuilder //条件集合
	Current *WhereCondition     //当前条件
	error   error
}

func NewEmptyConditionBuilder(or bool) *ConditionBuilder {
	return &ConditionBuilder{
		Or:      or,
		Items:   nil,
		Current: nil,
		error:   nil,
	}
}

func NewConditionBuilder(or bool, condition *WhereCondition) *ConditionBuilder {
	var builder = &ConditionBuilder{
		Or:      or,
		Items:   nil,
		Current: condition,
		error:   nil,
	}

	if condition == nil {
		return builder.Error("newConditionBuilder condition is nil")
	}

	return builder
}

func (c *ConditionBuilder) SetCondition(condition *WhereCondition) *ConditionBuilder {
	c.Current = condition
	return c
}

// AddChildrenBuilder 添加子条件
func (c *ConditionBuilder) AddChildrenBuilder(builders ...*ConditionBuilder) *ConditionBuilder {
	if len(builders) == 0 {
		return c.Error("AddChildrenBuilder builders is empty")
	}

	c.Items = append(c.Items, builders...)
	return c
}

// AddChildrenCondition 添加子条件
func (c *ConditionBuilder) AddChildrenCondition(conditions ...*WhereCondition) *ConditionBuilder {
	if len(conditions) == 0 {
		return c.Error("AddChildrenBuilder conditions is empty")
	}

	for _, condition := range conditions {
		c.AddChildrenBuilder(&ConditionBuilder{Current: condition})
	}

	return c
}

// BuildSql 生成sql
func (c *ConditionBuilder) BuildSql() (string, []interface{}, error) {
	if c == nil {
		return "", nil, errors.New("没有任何条件")
	}

	var compareSymbols = ""
	if c.Or {
		compareSymbols = "OR "
	} else {
		compareSymbols = "AND "
	}

	//没有子项，条件就是本身；有子项则用子项
	if len(c.Items) == 0 {
		if c.Current == nil {
			return "", nil, errors.New("没有任何条件")
		}

		return (*c.Current).GetSql(), (*c.Current).GetParams(), nil
	}

	var _sql = ""
	var _param = make([]interface{}, 0)
	for _, item := range c.Items {
		sql, param, err := item.BuildSql()
		if err != nil {
			return "", nil, err
		}

		_sql += compareSymbols + sql
		_param = append(_param, param...)
	}

	return "(" + _sql + ")", _param, nil
}

func (c *ConditionBuilder) Error(error string) *ConditionBuilder {
	c.error = errors.New(error)
	return c
}
