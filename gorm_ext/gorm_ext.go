package gorm_ext

import (
	"errors"
	"fmt"
	"github.com/ad313/go_ext/ext"
	"gorm.io/gorm/schema"
	"reflect"
	"sync"
)

type BuildOrmModelResult[T interface{}] struct {
	T     *T
	Error error
}

// 缓存实体对象，主要给NewQuery方法返回使用
var cache sync.Map

// columnMap 数据库表字段缓存
var columnMap = make(map[uintptr]string)

// BuildOrmModel 获取
func BuildOrmModel[T interface{}]() *BuildOrmModelResult[T] {
	modelTypeStr := reflect.TypeOf((*T)(nil)).Elem().String()
	if model, ok := cache.Load(modelTypeStr); ok {
		m, isReal := model.(*T)
		if isReal {
			return &BuildOrmModelResult[T]{T: m}
		}
	}

	t, _, ok := ext.IsType[T, schema.Tabler]()
	if ok == false {
		return &BuildOrmModelResult[T]{Error: errors.New("传入类型必须是实现了 TableName 的表实体")}
	}

	cache.Store(modelTypeStr, t)
	var cm = getColumnNameMap(t)
	for key, v := range cm {
		columnMap[key] = v
	}

	return &BuildOrmModelResult[T]{T: t}
}

// GetTableColumn 通过模型字段获取数据库字段
func GetTableColumn(column any) string {
	var v = reflect.ValueOf(column)
	var addr uintptr
	if v.Kind() == reflect.Pointer {
		addr = v.Pointer()
		n, ok := columnMap[addr]
		if ok {
			return n
		}
	} else {
		fmt.Println("column must be of type Pointer")
		return ""
	}

	return ""
}

func getColumnNameMap(model any) map[uintptr]string {
	var columnNameMap = make(map[uintptr]string)
	valueOf := reflect.ValueOf(model).Elem()
	typeOf := reflect.TypeOf(model).Elem()
	for i := 0; i < valueOf.NumField(); i++ {
		field := typeOf.Field(i)
		// 如果当前实体嵌入了其他实体，同样需要缓存它的字段名
		if field.Anonymous {
			// 如果存在多重嵌套，通过递归方式获取他们的字段名
			subFieldMap := getSubFieldColumnNameMap(valueOf, field)
			for pointer, columnName := range subFieldMap {
				columnNameMap[pointer] = columnName
			}
		} else {
			// 获取对象字段指针值
			pointer := valueOf.Field(i).Addr().Pointer()
			columnName := parseColumnName(field)
			if columnName != "" {
				columnNameMap[pointer] = columnName
			}
		}
	}
	return columnNameMap
}

// 递归获取嵌套字段名
func getSubFieldColumnNameMap(valueOf reflect.Value, field reflect.StructField) map[uintptr]string {
	result := make(map[uintptr]string)
	modelType := field.Type
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	for j := 0; j < modelType.NumField(); j++ {
		subField := modelType.Field(j)
		if subField.Anonymous {
			nestedFields := getSubFieldColumnNameMap(valueOf, subField)
			for key, value := range nestedFields {
				result[key] = value
			}
		} else {
			pointer := valueOf.FieldByName(modelType.Field(j).Name).Addr().Pointer()
			name := parseColumnName(modelType.Field(j))
			result[pointer] = name
		}
	}

	return result
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
func getSqlSm() string {
	//todo
	return ""
	//switch config.CFG.DB.Db.Type {
	//case "mysql":
	//	return "'"
	////case "clickhouse":
	////	instance.DB = ch.NewClickHouse(cfg)
	////	break
	//case "sqlite":
	//	return "'"
	//case "dm":
	//	return "\""
	////case "postgres", "pgsql":
	////	instance.DB = postgres.NewPostgres(cfg)
	//default:
	//
	//	break
	//}
	//
	//return "'"
}
