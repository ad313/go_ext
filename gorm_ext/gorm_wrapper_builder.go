package gorm_ext

import (
	"context"
	"errors"
	"fmt"
	"github.com/ad313/go_ext/ext"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

// PagerList 分页数据结果模型
type PagerList[T interface{}] struct {
	Page       int32  `json:"page" form:"page"`               //页码
	PageSize   int32  `json:"page_size" form:"page_size"`     //分页条数
	TotalCount int32  `json:"total_count" form:"total_count"` //总条数
	Order      string `json:"order" form:"order"`             //排序字段
	Data       *[]*T  `json:"data" form:"data"`               //数据项
}

// Pager 分页数据请求模型
type Pager struct {
	Page     int32  `json:"page" form:"page"`           //页码
	PageSize int32  `json:"page_size" form:"page_size"` //分页条数
	Order    string `json:"order" form:"order"`         //排序字段
	Keyword  string `json:"keyword" form:"keyword"`     //关键词
}

type ormWrapperBuilder[T interface{}] struct {
	wrapper *OrmWrapper[T]

	TableName  string //表名
	TableAlias string //表别名

	ctx       context.Context
	DbContext *gorm.DB
	isOuterDb bool //是否外部传入db

	where          [][]any          //普通条件
	WhereCondition []WhereCondition //condition条件
	leftJoin       []*leftJoinModel //leftJoin 集合
	selectColumns  []string         //select 字段集合
	groupByColumns []string         //group by 字段集合
	orderByColumns []string         //order by 字段集合

	isUnscoped bool //和gorm一样，忽略软删除字段
}

// leftJoinModel 左连接条件，会自动加上软删除字段
type leftJoinModel struct {
	Table     schema.Tabler //连接表，右表
	tableName string        //右表表名
	Alias     string        //右表表别名
	Left      string        //左表字段
	Right     string        //右表字段
	ext       string        //扩展字段，比如有软删除字段，这里加上软删除sql
}

func (a *ormWrapperBuilder[T]) addWhere(query interface{}, args []interface{}) {
	if a.where == nil {
		a.where = make([][]interface{}, 0)
	}

	if query == nil {
		a.wrapper.Error = errors.New("query 条件不能为空")
		return
	}

	a.where = append(a.where, append([]interface{}{query}, args...))
}

func (a *ormWrapperBuilder[T]) addWhereWithWhereCondition(condition WhereCondition) {
	if a.WhereCondition == nil {
		a.WhereCondition = make([]WhereCondition, 0)
	}
	a.WhereCondition = append(a.WhereCondition, condition)

	//if a.where == nil {
	//	a.where = make([][]interface{}, 0)
	//}
	//
	//sql, param, err := condition.BuildSql(_dbType)
	//if err != nil {
	//	a.wrapper.Error = errors.New("query 条件不能为空")
	//	return
	//}
	//
	//a.addWhere(sql, param)
}

func (a *ormWrapperBuilder[T]) mergeColumnName(column string, columnAlias string, tableAlias string) string {
	if tableAlias != "" {
		column = formatSqlName(tableAlias, _dbType) + "." + column
	}

	if columnAlias != "" {
		column += " as " + getSqlSm(_dbType) + columnAlias + getSqlSm(_dbType)
	}

	return column
}

func (a *ormWrapperBuilder[T]) mergeColumnNameWithFunc(column string, columnAlias string, tableAlias string, f string) string {
	var sql, _ = mergeTableColumnWithFunc(column, tableAlias, f, _dbType)
	if columnAlias != "" {
		sql += " as " + getSqlSm(_dbType) + columnAlias + getSqlSm(_dbType)
	}

	return column
}

// 设置主表，针对没有主动设置表别名，这里自动加上表名称做表别名
func (a *ormWrapperBuilder[T]) buildModel() {
	//没有手动设置表别名，这里判断是否需要加：left join、exists
	if a.TableAlias == "" {
		//leftJoin
		if len(a.leftJoin) > 0 {
			a.TableAlias = a.TableName
		} else {
			//exists
			if len(a.WhereCondition) > 0 {
				for _, condition := range a.WhereCondition {
					_, ok := condition.(*ExistsCondition)
					if ok {
						a.TableAlias = a.TableName
					}
				}
			}
		}
	}

	if a.TableAlias != "" {
		a.DbContext = a.DbContext.Model(new(T)).Table(a.TableName + " as " + a.TableAlias)
	} else {
		a.DbContext = a.DbContext.Model(new(T))
	}
}

func (a *ormWrapperBuilder[T]) buildWhere() {
	if a.where == nil {
		a.where = make([][]interface{}, 0)
	}

	if len(a.WhereCondition) > 0 {
		for _, condition := range a.WhereCondition {
			sql, param, err := condition.BuildSql(_dbType, a.isUnscoped)
			if err != nil {
				a.wrapper.Error = errors.New("query 条件不能为空")
				return
			}

			a.addWhere(sql, param)
		}
	}

	for _, items := range a.where {
		if len(items) == 0 {
			continue
		}

		if len(items) == 1 {
			a.DbContext = a.DbContext.Where(items[0])
		} else {
			a.DbContext = a.DbContext.Where(items[0], items[1:]...)
		}
	}
}

func (a *ormWrapperBuilder[T]) buildLeftJoin() {
	if len(a.leftJoin) > 0 {
		for _, join := range a.leftJoin {
			a.DbContext = a.DbContext.
				Joins(fmt.Sprintf("left join %v as %v on %v = %v%v",
					formatSqlName(join.tableName, _dbType),
					formatSqlName(join.Alias, _dbType),
					join.Left,
					join.Right,
					ext.ChooseTrueValue(a.isUnscoped, "", join.ext)))
		}

		a.DbContext = a.DbContext.Distinct()
	}
}

func (a *ormWrapperBuilder[T]) buildSelect() {
	if a.selectColumns != nil && len(a.selectColumns) > 0 {
		a.DbContext = a.DbContext.Select(strings.Join(a.selectColumns, ","))
	}
}

func (a *ormWrapperBuilder[T]) buildOrderBy() {
	if a.orderByColumns != nil && len(a.orderByColumns) > 0 {
		a.DbContext = a.DbContext.Order(strings.Join(a.orderByColumns, ","))
	}
}

func (a *ormWrapperBuilder[T]) buildGroupBy() {
	if a.groupByColumns != nil && len(a.groupByColumns) > 0 {
		//特殊处理一个参数的情况，否则报错
		if len(a.groupByColumns) == 1 {
			a.DbContext = a.DbContext.Group(strings.ReplaceAll(a.groupByColumns[0], getSqlSm(_dbType), ""))
		} else {
			a.DbContext = a.DbContext.Group(strings.Join(a.groupByColumns, ","))
		}
	}
}

// Build 创建 gorm sql
func (a *ormWrapperBuilder[T]) Build() *gorm.DB {
	a.buildWhere()
	a.buildSelect()
	a.buildLeftJoin()
	a.buildOrderBy()
	a.buildGroupBy()
	return a.DbContext
}

// BuildForQuery 创建 gorm sql
func (a *ormWrapperBuilder[T]) BuildForQuery() *gorm.DB {
	a.buildModel()
	a.Build()
	return a.DbContext
}
