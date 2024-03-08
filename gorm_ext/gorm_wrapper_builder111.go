package gorm_ext

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"gorm.io/gorm"
//	"gorm.io/gorm/schema"
//	"strings"
//)
//
//// PagerList 分页数据结果模型
//type PagerList[T interface{}] struct {
//	Page       int32  `json:"page" form:"page"`               //页码
//	PageSize   int32  `json:"page_size" form:"page_size"`     //分页条数
//	TotalCount int32  `json:"total_count" form:"total_count"` //总条数
//	Order      string `json:"order" form:"order"`             //排序字段
//	Data       *[]*T  `json:"data" form:"data"`               //数据项
//	AnyData    any    `json:"any_data" form:"any_data"`       //数据项
//}
//
//// Pager 分页数据请求模型
//type Pager struct {
//	Page     int32  `json:"page" form:"page"`           //页码
//	PageSize int32  `json:"page_size" form:"page_size"` //分页条数
//	Order    string `json:"order" form:"order"`         //排序字段
//	Keyword  string `json:"keyword" form:"keyword"`     //关键词
//}
//
//// OrmWrapper gorm包装器
//type OrmWrapper[T any] struct {
//	Error   error
//	Model   *T
//	builder *OrmWrapperBuilder[T]
//}
//
//type leftJoinModel struct {
//	Table     schema.Tabler
//	TableName string
//	Alias     string
//	Left      string
//	Right     string
//}
//
//type ExistsModel struct {
//	Table       *schema.Tabler
//	Columns     []*ExistsColumn
//	IsNotExists bool
//}
//
//func (a *ExistsModel) Set(table schema.Tabler, columns ...*ExistsColumn) *ExistsModel {
//	a.Table = &table
//	a.Columns = columns
//	return a
//}
//
//type ExistsColumn struct {
//	InnerColumn    any         //exists 对应的表字段
//	OuterColumn    any         //外部表字段
//	OuterAlias     string      //外部表表别名
//	OuterValue     interface{} //外部直接传值
//	CompareSymbols string      //比较符号
//}
//
//type OrColumn struct {
//	Column         any
//	CompareSymbols string
//	Arg            interface{}
//	TableAlias     string
//	ExistsModel    *ExistsModel //exists 条件
//}
//
//type OrmWrapperBuilder[T interface{}] struct {
//	wrapper *OrmWrapper[T]
//
//	TableName string
//	DbContext *gorm.DB
//	isOuterDb bool //是否外部传入db
//	ctx       context.Context
//
//	where          [][]any
//	leftJoin       []*leftJoinModel
//	existsModels   []*ExistsModel
//	selectColumns  []string
//	groupByColumns []string
//	orderByColumns []string
//	orColumns      [][]*OrColumn
//}
//
//func (a *OrmWrapperBuilder[T]) buildModel() {
//	a.DbContext = a.DbContext.Model(new(T))
//}
//
//func (a *OrmWrapperBuilder[T]) buildWhere() {
//	//or
//	if a.orColumns != nil && len(a.orColumns) > 0 {
//		for _, columns := range a.orColumns {
//
//			var sql = ""
//			var args = make([]interface{}, 0)
//			for i, column := range columns {
//
//				var currentSql = ""
//				var currentArg = make([]interface{}, 0)
//
//				//exists
//				if column.ExistsModel != nil {
//					currentSql, currentArg = a.buildExistsMethod(column.ExistsModel)
//					if currentSql == "" {
//						continue
//					}
//					args = append(args, currentArg...)
//				} else {
//					if column.CompareSymbols == "" {
//						break
//					}
//
//					currentSql = a.mergeWhereString(column.Column, column.CompareSymbols, column.TableAlias)
//					args = append(args, a.mergeWhereValue(column.CompareSymbols, column.Arg))
//				}
//
//				sql += currentSql
//				if i < len(columns)-1 {
//					sql += " OR "
//				} else {
//					a.addWhere(sql, args)
//				}
//			}
//		}
//	}
//
//	for _, items := range a.where {
//		if len(items) == 0 {
//			continue
//		}
//
//		if len(items) == 1 {
//			a.DbContext = a.DbContext.Where(items[0])
//		} else {
//			a.DbContext = a.DbContext.Where(items[0], items[1:]...)
//		}
//	}
//}
//
//func (a *OrmWrapperBuilder[T]) buildLeftJoin() {
//	if a.leftJoin != nil && len(a.leftJoin) > 0 {
//		for _, join := range a.leftJoin {
//			//if i == 0 {
//			//	a.DbContext = a.DbContext.Table(formatTableAlias(a.TableName) + " as " + formatTableAlias(a.TableName)).
//			//		Joins(fmt.Sprintf("left join %v as %v on %v = %v", formatTableAlias(join.TableName), formatTableAlias(join.Alias), join.Left, join.Right))
//			//} else {
//			//	a.DbContext = a.DbContext.
//			//		Joins(fmt.Sprintf("left join %v as %v on %v = %v", formatTableAlias(join.TableName), formatTableAlias(join.Alias), join.Left, join.Right))
//			//}
//
//			a.DbContext = a.DbContext.
//				Joins(fmt.Sprintf("left join %v as %v on %v = %v", formatTableAlias(join.TableName), formatTableAlias(join.Alias), join.Left, join.Right))
//		}
//
//		a.DbContext = a.DbContext.Distinct()
//	}
//}
//
//func (a *OrmWrapperBuilder[T]) buildExists() {
//	if len(a.existsModels) == 0 {
//		return
//	}
//
//	for _, existsModel := range a.existsModels {
//		if len(existsModel.Columns) == 0 {
//			continue
//		}
//
//		sql, appendValues := a.buildExistsMethod(existsModel)
//
//		a.DbContext = a.DbContext.Where(sql, appendValues...)
//	}
//}
//
//func (a *OrmWrapperBuilder[T]) buildExistsMethod(existsModel *ExistsModel) (string, []interface{}) {
//	if existsModel == nil {
//		return "", nil
//	}
//
//	if len(existsModel.Columns) == 0 {
//		return "", nil
//	}
//
//	var sql = fmt.Sprintf("SELECT 1 FROM %v WHERE 1=1", formatTableAlias((*existsModel.Table).TableName()))
//	var appendValues = make([]interface{}, 0)
//	for _, column := range existsModel.Columns {
//
//		//外部传值的方式
//		if column.OuterValue != nil {
//			if column.CompareSymbols == "" {
//				column.CompareSymbols = Eq
//			}
//
//			sql += " AND " + a.mergeWhereString(column.InnerColumn, column.CompareSymbols)
//			appendValues = append(appendValues, a.mergeWhereValue(column.CompareSymbols, column.OuterValue))
//			continue
//		}
//
//		var innerColumn = a.resolveColumnName(column.InnerColumn)
//		if innerColumn == "" {
//			break
//		}
//
//		var outerColumn = a.resolveColumnName(column.OuterColumn)
//		if outerColumn == "" {
//			break
//		}
//
//		if column.OuterAlias != "" {
//			outerColumn = formatTableAlias(column.OuterAlias) + "." + outerColumn
//		}
//
//		sql += " AND " + innerColumn + " = " + outerColumn
//	}
//
//	////leftJoin 会指定主表；如果没有 leftJoin，此处指定主表
//	//if len(a.leftJoin) == 0 {
//	//	a.DbContext = a.DbContext.Table(formatTableAlias(a.TableName) + " as " + formatTableAlias(a.TableName))
//	//}
//
//	var first = "Exists"
//	if existsModel.IsNotExists {
//		first = "Not " + first
//	}
//	return fmt.Sprintf("%v (%v)", first, sql), appendValues
//}
//
//// 设置主表
//func (a *OrmWrapperBuilder[T]) setMainTable() {
//	var set = false
//	if len(a.leftJoin) > 0 || len(a.existsModels) > 0 {
//		set = true
//	} else if len(a.orColumns) > 0 {
//		for _, columns := range a.orColumns {
//			for _, column := range columns {
//				if column.ExistsModel != nil {
//					set = true
//					break
//				}
//			}
//
//			if set {
//				break
//			}
//		}
//	}
//
//	if set {
//		a.DbContext = a.DbContext.Table(formatTableAlias(a.TableName) + " as " + formatTableAlias(a.TableName))
//	}
//}
//
//func (a *OrmWrapperBuilder[T]) buildSelect() {
//	if a.selectColumns != nil && len(a.selectColumns) > 0 {
//		a.DbContext = a.DbContext.Select(strings.Join(a.selectColumns, ","))
//	}
//}
//
//func (a *OrmWrapperBuilder[T]) buildOrderBy() {
//	if a.orderByColumns != nil && len(a.orderByColumns) > 0 {
//		a.DbContext = a.DbContext.Order(strings.Join(a.orderByColumns, ","))
//	}
//}
//
//func (a *OrmWrapperBuilder[T]) buildGroupBy() {
//	if a.groupByColumns != nil && len(a.groupByColumns) > 0 {
//		//特殊处理一个参数的情况，否则报错
//		if len(a.groupByColumns) == 1 {
//			a.DbContext = a.DbContext.Group(strings.ReplaceAll(a.groupByColumns[0], getSqlSm(), ""))
//		} else {
//			a.DbContext = a.DbContext.Group(strings.Join(a.groupByColumns, ","))
//		}
//	}
//}
//
//func (a *OrmWrapperBuilder[T]) addWhere(query interface{}, args []interface{}) {
//	if a.where == nil {
//		a.where = make([][]interface{}, 0)
//	}
//
//	if query != nil {
//		a.where = append(a.where, append([]interface{}{query}, args...))
//	}
//}
//
//// mergeWhereString 组合 where 条件
//func (a *OrmWrapperBuilder[T]) mergeWhereString(column any, compareSymbols string, tableAlias ...string) string {
//	var name = a.resolveColumnName(column)
//	var valueExpress = "?"
//	switch compareSymbols {
//	case "IN":
//		valueExpress = "(?)"
//		break
//	case "NOT IN":
//		valueExpress = "(?)"
//		break
//	case "IS NULL":
//		valueExpress = ""
//		break
//	case "IS NOT NULL":
//		valueExpress = ""
//		break
//	}
//
//	var table = ""
//	if tableAlias != nil && len(tableAlias) > 0 && tableAlias[0] != "" {
//		table = formatTableAlias(tableAlias[0]) + "."
//	}
//
//	return fmt.Sprintf("%v%v %v %v", table, name, compareSymbols, valueExpress)
//}
//
//// mergeWhereValue 处理查询值，当 IN 条件时，拼接 %
//func (a *OrmWrapperBuilder[T]) mergeWhereValue(compareSymbols string, value interface{}) interface{} {
//	if value == nil {
//		return value
//	}
//	v, ok := value.(string)
//	if ok == false {
//		return value
//	}
//
//	switch compareSymbols {
//	case "LIKE":
//		return "%" + v + "%"
//	case "NOT LIKE":
//		return "%" + v + "%"
//	case "STARTWITH":
//		return v + "%"
//	case "EndWith":
//		return "%" + v
//	}
//
//	return value
//}
//
//func (a *OrmWrapperBuilder[T]) mergeColumnName(column string, columnAlias string, tableAlias string) string {
//	if tableAlias != "" {
//		column = formatTableAlias(tableAlias) + "." + column
//	}
//
//	if columnAlias != "" {
//		column += " as " + getSqlSm() + columnAlias + getSqlSm()
//	}
//
//	return column
//}
//
//// mergeColumnAndTable 组合表名和字段名
//func (a *OrmWrapperBuilder[T]) mergeColumnAndTable(column any, tableAlias ...string) string {
//	var name = a.resolveColumnName(column)
//
//	if tableAlias != nil && len(tableAlias) > 0 {
//		name = formatTableAlias(tableAlias[0]) + "." + name
//	}
//
//	return name
//}
//
//// resolveColumnName 从缓存获取数据库字段名称：如果不能匹配，则返回 string 值
//func (a *OrmWrapperBuilder[T]) resolveColumnName(column any) string {
//	var name = ResolveTableColumnName(column)
//	if name == "" {
//		a.wrapper.Error = errors.New("未获取到字段名称")
//		return ""
//	}
//	return name
//}
//
//// clear 清除所有条件
//func (a *OrmWrapperBuilder[T]) clear() {
//	if a == nil {
//		return
//	}
//
//	a.wrapper.Error = nil
//	a.orderByColumns = make([]string, 0)
//	a.selectColumns = make([]string, 0)
//	a.groupByColumns = make([]string, 0)
//	a.leftJoin = make([]*leftJoinModel, 0)
//	a.existsModels = make([]*ExistsModel, 0)
//	a.where = make([][]any, 0)
//}
//
////// ToPagerFromProto 把 proto 分页模型转换成本地分页模型
////func ToPagerFromProto(source any) *Pager {
////	var pager = &Pager{
////		Page:     1,
////		PageSize: 20,
////	}
////
////	if source == nil {
////		return pager
////	}
////
////	err := mapper.MapProtoToStruct(source, &pager)
////	if err != nil {
////		return pager
////	}
////
////	if pager.Page <= 0 {
////		pager.Page = 1
////	}
////
////	if pager.PageSize <= 0 {
////		pager.PageSize = 20
////	}
////
////	return pager
////}
////
////// ToPagerFromStruct 把 struct 分页模型转换成本地分页模型
////func ToPagerFromStruct(source any) *Pager {
////	var pager = &Pager{
////		Page:     1,
////		PageSize: 20,
////	}
////
////	if source == nil {
////		return pager
////	}
////
////	err := mapper.MapTo(source, &pager)
////	if err != nil {
////		return pager
////	}
////
////	if pager.Page <= 0 {
////		pager.Page = 1
////	}
////
////	if pager.PageSize <= 0 {
////		pager.PageSize = 20
////	}
////
////	return pager
////}
