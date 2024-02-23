package gorm_ext

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

var _db *gorm.DB = nil

// Init 初始化塞入db
func Init(db *gorm.DB) {
	_db = db
}

// BuildOrmWrapper 创建gorm包装器
func BuildOrmWrapper[T any](ctx context.Context, db ...*gorm.DB) *OrmWrapper[T] {
	var wrapper = &OrmWrapper[T]{}

	//创建模型
	var buildResult = BuildOrmModel[T]()
	wrapper.Model = buildResult.T
	wrapper.Error = buildResult.Error
	wrapper.builder = &OrmWrapperBuilder[T]{
		wrapper:        wrapper,
		where:          make([][]any, 0),
		leftJoin:       make([]*leftJoinModel, 0),
		selectColumns:  make([]string, 0),
		groupByColumns: make([]string, 0),
		orderByColumns: make([]string, 0)}

	wrapper.SetDbContext(ctx, db...)

	if wrapper.Error == nil {
		model, ok := IsTypeByValue[schema.Tabler](*(wrapper.Model))
		if ok {
			wrapper.builder.TableName = (*model).TableName()
		} else {
			wrapper.Error = errors.New("传入类型必须是实现了 TableName 的表实体")
		}
	}

	return wrapper
}

func (a *OrmWrapper[T]) SetDbContext(ctx context.Context, db ...*gorm.DB) *OrmWrapper[T] {
	if len(db) > 0 {
		a.builder.DbContext = db[0]
		a.builder.isOuterDb = true
	}

	a.builder.ctx = ctx

	return a
}

func (a *OrmWrapper[T]) GetNewDbContext(ctx context.Context) *gorm.DB {
	return _db.WithContext(ctx)
}

// Where gorm 原生查询
func (a *OrmWrapper[T]) Where(query interface{}, args ...interface{}) *OrmWrapper[T] {
	a.builder.addWhere(query, args)
	return a
}

// WhereIf gorm 原生查询，加入 bool 条件控制
func (a *OrmWrapper[T]) WhereIf(do bool, query interface{}, args ...interface{}) *OrmWrapper[T] {
	if do && query != nil {
		return a.Where(query, args...)
	}
	return a
}

// WhereIfNotNil gorm 原生查询，值为 nil 时跳过
func (a *OrmWrapper[T]) WhereIfNotNil(query interface{}, arg interface{}) *OrmWrapper[T] {
	if arg != nil && query != nil {
		return a.Where(query, arg)
	}

	return a
}

// WhereByColumn 通过字段查询，连表时支持传入表别名
func (a *OrmWrapper[T]) WhereByColumn(column any, compareSymbols string, arg interface{}, tableAlias ...string) *OrmWrapper[T] {
	if arg != nil && column != nil && compareSymbols != "" {
		a.builder.addWhere(a.builder.mergeWhereString(column, compareSymbols, tableAlias...), []interface{}{a.builder.mergeWhereValue(compareSymbols, arg)})
	}
	return a
}

// WhereByColumnIf 通过字段查询，连表时支持传入表别名
func (a *OrmWrapper[T]) WhereByColumnIf(do bool, column any, compareSymbols string, arg interface{}, tableAlias ...string) *OrmWrapper[T] {
	if do {
		a.WhereByColumn(column, compareSymbols, arg, tableAlias...)
	}
	return a
}

// OrColumnIf Or条件，外部 and，内部 or
func (a *OrmWrapper[T]) OrColumnIf(do bool, columns ...*OrColumn) *OrmWrapper[T] {
	if do {
		return a.OrColumn(columns...)
	}

	return a
}

// OrColumn Or条件，外部 and，内部 or
func (a *OrmWrapper[T]) OrColumn(columns ...*OrColumn) *OrmWrapper[T] {
	if len(columns) >= 2 {
		if a.builder.orColumns == nil {
			a.builder.orColumns = make([][]*OrColumn, 0)
		}

		a.builder.orColumns = append(a.builder.orColumns, columns)
	}

	return a
}

// LeftJoin 左连表
func (a *OrmWrapper[T]) LeftJoin(table schema.Tabler, alias string, leftColumn any, rightColumn any, selectColumns ...any) *OrmWrapper[T] {
	if a.builder.leftJoin == nil {
		a.builder.leftJoin = make([]*leftJoinModel, 0)
	}

	if table == nil || leftColumn == nil || rightColumn == nil {
		return a
	}

	var left = a.builder.resolveColumnName(leftColumn)
	if left == "" {
		a.Error = errors.New("LeftJoin 未获取到左边字段")
		return a
	}

	var right = a.builder.resolveColumnName(rightColumn)
	if right == "" {
		a.Error = errors.New("LeftJoin 未获取到右边字段")
		return a
	}

	if alias == "" {
		alias = table.TableName()
	}

	var leftTableName string
	if len(a.builder.leftJoin) == 0 {
		leftTableName = a.builder.TableName
	} else {
		leftTableName = a.builder.leftJoin[len(a.builder.leftJoin)-1].Alias
	}

	var joinModel = &leftJoinModel{
		Table:     table,
		TableName: table.TableName(),
		Alias:     alias,
		Left:      formatTableAlias(leftTableName) + "." + left,
		Right:     formatTableAlias(alias) + "." + right,
	}
	a.builder.leftJoin = append(a.builder.leftJoin, joinModel)

	if selectColumns != nil && len(selectColumns) > 0 {
		for _, column := range selectColumns {
			var name = a.builder.resolveColumnName(column)
			if name != "" {
				a.builder.selectColumns = append(a.builder.selectColumns, formatTableAlias(joinModel.Alias)+"."+name)
			} else {
				a.Error = errors.New("LeftJoin 未获取到 select 字段")
			}
		}
	}

	return a
}

// LeftJoinIf 左连表
func (a *OrmWrapper[T]) LeftJoinIf(do bool, table schema.Tabler, alias string, leftColumn any, rightColumn any, selectColumns ...any) *OrmWrapper[T] {
	if do {
		a.LeftJoin(table, alias, leftColumn, rightColumn, selectColumns...)
	}

	return a
}

// Select 查询主表字段。如果要查询连表的字段，则在 LeftJoin 时传入
func (a *OrmWrapper[T]) Select(selectColumns ...interface{}) *OrmWrapper[T] {
	if selectColumns == nil || len(selectColumns) == 0 {
		return a
	}

	var isSampleQuery = a.builder.leftJoin == nil || len(a.builder.leftJoin) == 0
	var table = ""
	if isSampleQuery == false {
		table = a.builder.TableName
	}

	return a.SelectWithTableAlias(table, selectColumns...)
}

// SelectWithTableAlias 传入表别名，查询此表下的字段
func (a *OrmWrapper[T]) SelectWithTableAlias(tableAlias string, selectColumns ...interface{}) *OrmWrapper[T] {
	if selectColumns == nil || len(selectColumns) == 0 {
		return a
	}

	for _, column := range selectColumns {
		var name = a.builder.resolveColumnName(column)
		if name == "" {
			a.Error = errors.New("未获取到字段名称")
			continue
		}

		a.builder.selectColumns = append(a.builder.selectColumns, a.builder.mergeColumnName(name, "", tableAlias))
	}

	return a
}

// SelectColumn 单次查询一个字段，可传入 字段别名，表名；可多次调用
func (a *OrmWrapper[T]) SelectColumn(selectColumn any, columnAlias string, tableAlias string) *OrmWrapper[T] {
	var name = a.builder.resolveColumnName(selectColumn)
	if name == "" {
		a.Error = errors.New("未获取到字段名称")
		return a
	}

	a.builder.selectColumns = append(a.builder.selectColumns, a.builder.mergeColumnName(name, columnAlias, tableAlias))

	return a
}

// SelectColumnOriginal 单次查询一个字段，可传入 字段别名，表名；可多次调用；不处理字段名
func (a *OrmWrapper[T]) SelectColumnOriginal(selectColumn string, columnAlias string, tableAlias string) *OrmWrapper[T] {
	if selectColumn == "" {
		a.Error = errors.New("未获取到字段名称")
		return a
	}

	a.builder.selectColumns = append(a.builder.selectColumns, a.builder.mergeColumnName(selectColumn, columnAlias, tableAlias))

	return a
}

// GroupBy 可多次调用，按照调用顺序排列字段
func (a *OrmWrapper[T]) GroupBy(column any, tableAlias ...string) *OrmWrapper[T] {
	if a.builder.groupByColumns == nil {
		a.builder.groupByColumns = make([]string, 0)
	}

	var name = a.builder.mergeColumnAndTable(column, tableAlias...)
	if name != "" {
		a.builder.groupByColumns = append(a.builder.groupByColumns, name)
	} else {
		a.Error = errors.New("未获取到 GroupBy 字段名称")
	}

	return a
}

// OrderBy 可多次调用，按照调用顺序排列字段
func (a *OrmWrapper[T]) OrderBy(column any, tableAlias ...string) *OrmWrapper[T] {
	if a.builder.orderByColumns == nil {
		a.builder.orderByColumns = make([]string, 0)
	}

	if column == nil {
		return a
	}

	var name = a.builder.mergeColumnAndTable(column, tableAlias...)
	if name != "" {
		a.builder.orderByColumns = append(a.builder.orderByColumns, name)
	} else {
		a.Error = errors.New("未获取到 OrderBy 字段名称")
	}

	return a
}

// OrderByDesc 可多次调用，按照调用顺序排列字段
func (a *OrmWrapper[T]) OrderByDesc(column any, tableAlias ...string) *OrmWrapper[T] {
	if a.builder.orderByColumns == nil {
		a.builder.orderByColumns = make([]string, 0)
	}

	var name = a.builder.mergeColumnAndTable(column, tableAlias...)
	if name != "" {
		a.builder.orderByColumns = append(a.builder.orderByColumns, name+" desc")
	} else {
		a.Error = errors.New("未获取到 OrderByDesc 字段名称")
	}

	return a
}

// WhereExists exists 语句
func (a *OrmWrapper[T]) WhereExists(table schema.Tabler, columns ...*ExistsColumn) *OrmWrapper[T] {
	if table == nil || len(columns) < 1 {
		return a
	}

	if a.builder.existsModels == nil {
		a.builder.existsModels = make([]*ExistsModel, 0)
	}

	a.builder.existsModels = append(a.builder.existsModels, &ExistsModel{
		Table:   &table,
		Columns: columns,
	})

	return a
}

// WhereNotExists exists 语句
func (a *OrmWrapper[T]) WhereNotExists(table schema.Tabler, columns ...*ExistsColumn) *OrmWrapper[T] {
	if table == nil || len(columns) < 1 {
		return a
	}

	if a.builder.existsModels == nil {
		a.builder.existsModels = make([]*ExistsModel, 0)
	}

	a.builder.existsModels = append(a.builder.existsModels, &ExistsModel{
		Table:       &table,
		Columns:     columns,
		IsNotExists: true,
	})

	return a
}

// WhereExistsIf exists 语句
func (a *OrmWrapper[T]) WhereExistsIf(do bool, table schema.Tabler, columns ...*ExistsColumn) *OrmWrapper[T] {
	if do {
		return a.WhereExists(table, columns...)
	}

	return a
}

// Count 查询总条数
func (a *OrmWrapper[T]) Count() (int64, error) {
	//创建语句过程中的错误
	if a.Error != nil {
		return 0, a.Error
	}

	//Build sql
	a.BuildForQuery()

	var result int64
	//First 会自动添加主键排序
	err := a.builder.DbContext.Count(&result).Error
	if err != nil {
		return 0, err
	}

	return result, nil
}

// FirstOrDefault 返回第一条，没命中返回nil
func (a *OrmWrapper[T]) FirstOrDefault() (*T, error) {

	//创建语句过程中的错误
	if a.Error != nil {
		return nil, a.Error
	}

	//Build sql
	a.BuildForQuery()

	var result = new(T)
	//First 会自动添加主键排序
	err := a.builder.DbContext.Take(result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}

// FirstOrDefaultCustom 返回第一条，没命中返回nil
func (a *OrmWrapper[T]) FirstOrDefaultCustom(result any) error {

	//创建语句过程中的错误
	if a.Error != nil {
		return a.Error
	}

	//Build sql
	a.BuildForQuery()

	//First 会自动添加主键排序
	err := a.builder.DbContext.Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return err
	}

	return nil
}

// ToList 返回列表
func (a *OrmWrapper[T]) ToList(scan ...func(db *gorm.DB) error) ([]*T, error) {

	//创建语句过程中的错误
	if a.Error != nil {
		return nil, a.Error
	}

	//Build sql
	a.BuildForQuery()

	if scan != nil && len(scan) > 0 {
		return nil, scan[0](a.builder.DbContext)
	}

	var list = make([]*T, 0)
	err := a.builder.DbContext.Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// ToPagerList 分页查询，返回当前实体的分页结果
func (a *OrmWrapper[T]) ToPagerList(pager *Pager) (*PagerList[T], error) {

	//创建语句过程中的错误
	if a.Error != nil {
		return nil, a.Error
	}

	if pager == nil {
		return nil, errors.New("传入分页数据不能为空")
	}

	//包含空格 asc desc
	if strings.Contains(pager.Order, " ") {
		var arr = strings.Split(pager.Order, " ")
		if strings.ToUpper(arr[1]) == "DESC" {
			a.OrderByDesc(arr[0])
		} else {
			a.OrderBy(arr[0])
		}
	} else {
		a.OrderBy(pager.Order)
	}

	//Build sql
	a.BuildForQuery()

	if pager.Page <= 0 {
		pager.Page = 1
	}

	if pager.PageSize <= 0 {
		pager.PageSize = 20
	}

	//总条数
	var total int64
	var err error

	//left join 加上 distinct
	if len(a.builder.leftJoin) > 0 {
		var query = a.builder.DbContext
		err = a.GetNewDbContext(a.builder.ctx).Table("(?) as leftJoinTableWrapper", query.Distinct()).Count(&total).Error
	} else {
		err = a.builder.DbContext.Count(&total).Error
	}
	if err != nil {
		return nil, err
	}

	var result = &PagerList[T]{
		Page:       pager.Page,
		PageSize:   pager.PageSize,
		TotalCount: int32(total),
		Order:      pager.Order,
	}

	var data = make([]*T, 0)
	result.Data = &data

	if result.TotalCount > 0 {
		err = a.builder.DbContext.Offset(int((pager.Page - 1) * pager.PageSize)).Limit(int(pager.PageSize)).Scan(&result.Data).Error
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ToPagerListCustom 分页查询，返回自定义实体的分页结果
func (a *OrmWrapper[T]) ToPagerListCustom(pager *Pager, scan func(db *gorm.DB) error) (*PagerList[T], error) {

	//创建语句过程中的错误
	if a.Error != nil {
		return nil, a.Error
	}

	if pager == nil {
		return nil, errors.New("传入分页数据不能为空")
	}

	a.OrderBy(pager.Order)

	//Build sql
	a.BuildForQuery()

	if pager.Page <= 0 {
		pager.Page = 1
	}

	if pager.PageSize <= 0 {
		pager.PageSize = 20
	}

	//总条数
	var total int64
	err := a.builder.DbContext.Count(&total).Error
	if err != nil {
		return nil, err
	}

	var result = &PagerList[T]{
		Page:       pager.Page,
		PageSize:   pager.PageSize,
		TotalCount: int32(total),
		Order:      pager.Order,
	}

	if scan == nil {
		return nil, errors.New("scan 函数不能为空")
	}

	if result.TotalCount > 0 {
		err = scan(a.builder.DbContext.Offset(int((pager.Page - 1) * pager.PageSize)).Limit(int(pager.PageSize)))
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (a *OrmWrapper[T]) Unscoped() *OrmWrapper[T] {
	a.builder.DbContext = a.builder.DbContext.Unscoped()
	return a
}

//// ToPagerListCustom 分页查询，返回自定义实体的分页结果
//func (a *OrmWrapper[T]) ToPagerListCustom(pager *Pager, data *any) (*PagerList[any], error) {
//
//	//创建语句过程中的错误
//	if a.Error != nil {
//		return nil, a.Error
//	}
//
//	if pager == nil {
//		return nil, errors.New("传入分页数据不能为空")
//	}
//
//	a.OrderBy(pager.Order)
//
//	//Build sql
//	a.BuildForQuery()
//
//	if pager.Page <= 0 {
//		pager.Page = 1
//	}
//
//	if pager.PageSize <= 0 {
//		pager.PageSize = 20
//	}
//
//	//总条数
//	var total int64
//	err := a.builder.DbContext.Count(&total).Error
//	if err != nil {
//		return nil, err
//	}
//
//	var result = &PagerList[any]{
//		Page:       pager.Page,
//		PageSize:   pager.PageSize,
//		TotalCount: int32(total),
//		Order:      pager.Order,
//	}
//
//	if result.TotalCount > 0 {
//		err = a.builder.DbContext.Offset(int((pager.Page - 1) * pager.PageSize)).Limit(int(pager.PageSize)).Scan(data).Error
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	result.AnyData = data
//
//	return result, nil
//}

// Update 更新，传了字段只更新出入字段，否则更新全部
func (a *OrmWrapper[T]) Update(item *T, updateColumns ...interface{}) (int64, error) {
	if item == nil {
		return 0, nil
	}

	var isUpdateAll = false
	if len(updateColumns) > 0 {
		a.Select(updateColumns...)
	} else {
		isUpdateAll = true
	}

	a.BuildForQuery()

	//创建语句过程中的错误
	if a.Error != nil {
		return 0, a.Error
	}

	var result *gorm.DB
	if isUpdateAll {
		result = a.builder.DbContext.Save(item)
		return result.RowsAffected, result.Error
	} else {
		result = a.builder.DbContext.UpdateColumns(item)
		return result.RowsAffected, result.Error
	}
}

// UpdateList 更新，传了字段只更新出入字段，否则更新全部
func (a *OrmWrapper[T]) UpdateList(items []*T, updateColumns ...interface{}) (int64, error) {
	if len(items) == 0 {
		return 0, nil
	}

	var total int64 = 0

	//外部开启了事务
	if a.builder.isOuterDb {
		for _, item := range items {
			c, err := a.Update(item, updateColumns...)
			if err != nil {
				return 0, err
			}

			total += c
		}

		return total, nil
	}

	//本地开事务
	//var dbContext =
	var db = a.builder.DbContext

	err := db.Transaction(func(tx *gorm.DB) error {
		for i, item := range items {
			//重新设置db
			if i == 0 {
				a.SetDbContext(a.builder.ctx, tx)
			}

			c, err := a.Update(item, updateColumns...)
			if err != nil {
				return err
			}

			total += c
		}

		a.SetDbContext(a.builder.ctx, db)

		return nil
	})

	if err != nil {
		return 0, err
	}

	return total, nil
}

// ResolveTableColumnName 从缓存获取数据库字段名称：如果不能匹配，则返回 string 值
func ResolveTableColumnName(column any) string {
	var kind = reflect.ValueOf(column).Kind()
	if kind == reflect.Pointer {
		var name = GetTableColumn(column)
		if name == "" {
			return ""
		}
		return getSqlSm() + name + getSqlSm()
	} else {
		if str, ok := column.(string); ok && str != "" {
			return getSqlSm() + str + getSqlSm()
		} else {
			return ""
		}
	}
}

// Build 创建 gorm sql
func (a *OrmWrapper[T]) Build() *gorm.DB {
	//清除条件
	//defer a.builder.clear()

	a.builder.setMainTable()
	a.builder.buildWhere()
	a.builder.buildSelect()
	a.builder.buildLeftJoin()
	a.builder.buildExists()
	a.builder.buildOrderBy()
	a.builder.buildGroupBy()
	return a.builder.DbContext
}

// BuildForQuery 创建 gorm sql
func (a *OrmWrapper[T]) BuildForQuery() *gorm.DB {
	a.builder.buildModel()
	a.Build()
	return a.builder.DbContext
}
