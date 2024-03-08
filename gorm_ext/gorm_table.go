package gorm_ext

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/ad313/go_ext/ext"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 表缓存
var tableCache sync.Map
var tableSchemaCache sync.Map

// 字段缓存
var columnCache = make(map[uintptr]string)

// TableInfo 表信息
type TableInfo[T interface{}] struct {
	T *T
	*TableSchema
}

type TableSchema struct {
	Name              string //表名
	DeletedColumnName string //软删除字段名
}

// GormTableResult 创建表缓存结果
type GormTableResult[T interface{}] struct {
	Table *TableInfo[T]
	Error error
}

// BuildGormTable 获取
func BuildGormTable[T interface{}]() *GormTableResult[T] {
	modelType := ext.GetPath[T]()
	if model, ok := tableCache.Load(modelType); ok {
		m, isReal := model.(*TableInfo[T])
		if isReal {
			return &GormTableResult[T]{Table: m}
		}
	}

	//验证 schema.Tabler
	t, table, ok := ext.IsType[T, schema.Tabler]()
	if !ok {
		return &GormTableResult[T]{Error: errors.New("T 必须实现 schema.Tabler")}
	}

	//处理字段
	m, softDeletedName := getColumnNameMap(t)
	for key, v := range m {
		columnCache[key] = v
	}

	//缓存
	var cache = &TableInfo[T]{
		T: t,
		TableSchema: &TableSchema{
			Name:              (*table).TableName(),
			DeletedColumnName: softDeletedName,
		},
	}
	tableCache.Store(modelType, cache)
	tableSchemaCache.Store(modelType, cache.TableSchema)

	return &GormTableResult[T]{Table: cache}
}

// GetTableColumn 通过模型字段获取数据库字段
func GetTableColumn(column any) string {
	var v = reflect.ValueOf(column)
	var addr uintptr
	if v.Kind() == reflect.Pointer {
		addr = v.Pointer()
		n, ok := columnCache[addr]
		if ok {
			return n
		}
	} else {
		log("column must be of type Pointer")
		return ""
	}

	return ""
}

// GetTableSchema 获取表元数据
func GetTableSchema(table schema.Tabler) *TableSchema {
	modelType := ext.GetPathByValue(table)
	if model, ok := tableSchemaCache.Load(modelType); ok {
		m, isReal := model.(*TableSchema)
		if isReal {
			return m
		}
	}

	return nil
}

func getColumnNameMap(model any) (map[uintptr]string, string) {
	var columnNameMap = make(map[uintptr]string)
	valueOf := reflect.ValueOf(model).Elem()
	typeOf := reflect.TypeOf(model).Elem()
	var softDeletedColumn = ""
	var childSoftDeletedColumn = ""

	for i := 0; i < valueOf.NumField(); i++ {
		field := typeOf.Field(i)
		// 如果当前实体嵌入了其他实体，同样需要缓存它的字段名
		if field.Anonymous {
			// 如果存在多重嵌套，通过递归方式获取他们的字段名
			subFieldMap, _childSoftDeletedColumn := getSubFieldColumnNameMap(valueOf, field)
			for pointer, columnName := range subFieldMap {
				columnNameMap[pointer] = columnName
			}

			if childSoftDeletedColumn == "" && _childSoftDeletedColumn != "" {
				childSoftDeletedColumn = _childSoftDeletedColumn
			}
		} else {
			// 获取对象字段指针值
			pointer := valueOf.Field(i).Addr().Pointer()
			columnName := parseColumnName(field)
			if columnName != "" {
				columnNameMap[pointer] = columnName
			}

			//判断软删除字段
			if isGormDeletedAt(field, valueOf) {
				softDeletedColumn = columnName
			}
		}
	}

	//优先用本级的软删除字段
	if softDeletedColumn == "" {
		softDeletedColumn = childSoftDeletedColumn
	}

	return columnNameMap, softDeletedColumn
}

func isGormDeletedAt(field reflect.StructField, valueOf reflect.Value) bool {
	//判断软删除字段
	_, ok := ext.IsTypeByValue[gorm.DeletedAt](valueOf.FieldByName(field.Name).Interface())
	return ok
}

// 递归获取嵌套字段名
func getSubFieldColumnNameMap(valueOf reflect.Value, field reflect.StructField) (map[uintptr]string, string) {
	result := make(map[uintptr]string)
	modelType := field.Type
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	var softDeletedColumn = ""
	var childSoftDeletedColumn = ""

	for j := 0; j < modelType.NumField(); j++ {
		subField := modelType.Field(j)
		if subField.Anonymous {
			nestedFields, _childSoftDeletedColumn := getSubFieldColumnNameMap(valueOf, subField)
			for key, value := range nestedFields {
				result[key] = value
			}

			if childSoftDeletedColumn == "" && _childSoftDeletedColumn != "" {
				childSoftDeletedColumn = _childSoftDeletedColumn
			}
		} else {
			pointer := valueOf.FieldByName(modelType.Field(j).Name).Addr().Pointer()
			name := parseColumnName(modelType.Field(j))
			result[pointer] = name

			//判断软删除字段
			if isGormDeletedAt(subField, valueOf) {
				softDeletedColumn = name
			}
		}
	}

	//优先用本级的软删除字段
	if softDeletedColumn == "" {
		softDeletedColumn = childSoftDeletedColumn
	}

	return result, softDeletedColumn
}

// 解析字段名称
func parseColumnName(field reflect.StructField) string {
	tagSetting := schema.ParseTagSetting(field.Tag.Get("gorm"), ";")
	name, ok := tagSetting["COLUMN"]
	if ok {
		return name
	}
	return ""
}

// getSqlSm 获取sql 中 数据库字段分隔符
func getSqlSm(dbType string) string {
	switch dbType {
	case MySql:
		return "`"
	case Sqlite:
		return "'"
	case Dm:
		return "\""
	case Postgres:
		return "'"
	case Sqlserver:
		return "'"
	default:
		break
	}

	return ""
}

// 记录日志
func log(content string) {
	fmt.Println(content)
}
