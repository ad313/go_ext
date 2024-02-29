package gorm_ext

import (
	"errors"
	"github.com/ad313/go_ext/ext"
	"strings"
)

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

func NewAndEmptyConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{
		Or:      false,
		Items:   nil,
		Current: nil,
		error:   nil,
	}
}

func NewOrEmptyConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{
		Or:      true,
		Items:   nil,
		Current: nil,
		error:   nil,
	}
}

// NewAndConditionBuilder 创建 Builder，当条件个数为1，则加在builder本身，大于1，则加在Items；内部关系是 And
func NewAndConditionBuilder(conditions ...*WhereCondition) *ConditionBuilder {
	return newConditionBuilder(false, conditions...)
}

// NewOrConditionBuilder 创建 Builder，当条件个数为1，则加在builder本身，大于1，则加在Items；内部关系是 Or
func NewOrConditionBuilder(conditions ...*WhereCondition) *ConditionBuilder {
	return newConditionBuilder(true, conditions...)
}

// SetCondition builder 设置本级条件
func (c *ConditionBuilder) SetCondition(condition *WhereCondition) *ConditionBuilder {
	c.Current = condition
	return c
}

// AddChildrenBuilder builder 添加子条件
func (c *ConditionBuilder) AddChildrenBuilder(builders ...*ConditionBuilder) *ConditionBuilder {
	if len(builders) == 0 {
		return c.Error("AddChildrenBuilder builders is empty")
	}

	c.Items = append(c.Items, builders...)
	return c
}

// AddChildrenCondition builder 添加子条件
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

	//没有子项，条件就是本身；有子项则用子项
	if len(c.Items) == 0 {
		if c.Current == nil {
			return "", nil, errors.New("没有任何条件")
		}

		return (*c.Current).GetSql(), (*c.Current).GetParams(), nil
	}

	var _sql = make([]string, 0)
	var _param = make([]interface{}, 0)
	for _, item := range c.Items {
		sql, param, err := item.BuildSql()
		if err != nil {
			return "", nil, err
		}

		_sql = append(_sql, sql)
		_param = append(_param, param...)
	}

	//条件符号
	var compareSymbols = ext.ChooseTrueValue(c.Or, " OR ", " AND ")

	//todo 一个条件时可省略括号
	return "(" + strings.Join(_sql, compareSymbols) + ")", _param, nil
}

func (c *ConditionBuilder) Error(error string) *ConditionBuilder {
	c.error = errors.New(error)
	return c
}

// 创建 Builder，当条件个数为1，则加在builder本身，大于1，则加在Items
func newConditionBuilder(or bool, conditions ...*WhereCondition) *ConditionBuilder {
	var builder = &ConditionBuilder{
		Or:    or,
		Items: nil,
		error: nil,
	}

	if len(conditions) == 0 {
		return builder.Error("newConditionBuilder conditions has no item")
	}

	if len(conditions) == 1 {
		return builder.SetCondition(conditions[0])
	}

	return builder.AddChildrenCondition(conditions...)
}
