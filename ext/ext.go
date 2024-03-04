package ext

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"sort"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// StringWrapper 字符串包装为可空类型
func StringWrapper(str string) *wrapperspb.StringValue {
	return &wrapperspb.StringValue{Value: str}
}

// StringValue 可空类型 转 string
func StringValue(str *wrapperspb.StringValue) string {
	if str == nil {
		return ""
	}
	return str.Value
}

// Int32Value 可空类型 转 string
func Int32Value(str *wrapperspb.Int32Value) int32 {
	if str == nil {
		return 0
	}
	return str.Value
}

// IsNullOrEmpty 判断字符串是空
func IsNullOrEmpty(str *wrapperspb.StringValue) bool {
	return str == nil || str.Value == ""
}

// IsNullOrEmptyString 判断字符串是空
func IsNullOrEmptyString(str string) bool {
	return str == ""
}

// IsNotNullOrEmpty 判断字符串非空
func IsNotNullOrEmpty(str *wrapperspb.StringValue) bool {
	return !IsNullOrEmpty(str)
}

func GetAuthorizationToken(ctx context.Context) string {
	return GetValueFromContext(ctx, "authorization")
}

func GetValueFromContext(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		var values = md.Get(key)
		if len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

// ValidatePassword 验证密码强度 同时包含大写字母、小写字母、数字、特殊字符，6-16位
func ValidatePassword(password string, min int, max int) bool {
	// 检查密码长度是否在6到16位之间
	if len(password) < min || len(password) > max {
		return false
	}

	// 检查密码是否包含大写字母、小写字母、数字和特殊字符
	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	hasSpecialChar := false

	for _, char := range password {
		switch char {
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			hasUppercase = true
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			hasLowercase = true
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			hasDigit = true
		default:
			matched, err := regexp.MatchString("^[^a-zA-Z0-9]*$", string(char))
			if err != nil {
				fmt.Println(err)
				return false
			}
			if matched {
				hasSpecialChar = true
			}
		}
	}

	// 如果密码同时包含大写字母、小写字母、数字和特殊字符，则返回true;否则返回false。
	return hasUppercase && hasLowercase && hasDigit && hasSpecialChar
}

// 判断是否是手机号码
func isPhoneNumber(s string) bool {
	pattern := `^1[3-9]\d{9}$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}

// 不包含特殊字符
func isValidString(s string) bool {
	pattern := `^[a-zA-Z0-9]+$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}

// GetValueByFieldName 获取指定字段的值
func GetValueByFieldName(o interface{}, name string) any {
	t := reflect.TypeOf(o)
	// 获取值
	v := reflect.ValueOf(o)
	// 可以获取所有属性
	// 获取结构体字段个数：t.NumField()
	for i := 0; i < t.NumField(); i++ {
		// 取每个字段
		f := t.Field(i)
		if f.Name == name {
			return v.Field(i).Interface()
		}
	}
	return nil
}

// ChooseTrueValue 模拟三元表达式，获取值
func ChooseTrueValue[T interface{}](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}

	return falseValue
}

// Distinct 去重
func Distinct[T comparable](slice []T) []T {
	if slice == nil {
		return make([]T, 0)
	}
	var m = make(map[T]int32)
	for _, e := range slice {
		_, ok := m[e]
		if ok {
			continue
		}
		m[e] = 0
	}

	var result = make([]T, 0)
	for t, _ := range m {
		result = append(result, t)
	}

	return result
}

// Contains 判断数组是否包含某个项
func Contains[T interface{}](slice []T, comparable func(t T) bool) bool {
	if slice == nil {
		return false
	}
	for _, e := range slice {
		if comparable(e) {
			return true
		}
	}

	return false
}

// ContainsByValue 判断数组是否包含某个项
func ContainsByValue[T comparable](slice []T, e T) bool {
	if slice == nil {
		return false
	}
	if Contains(slice, func(item T) bool { return item == e }) {
		return true
	}

	return false
}

// FindOne 从数组中返回第一个符合条件的项
func FindOne[T interface{}](slice []T, comparable func(t T) bool) T {
	if slice == nil {
		return *new(T)
	}
	for _, e := range slice {
		if comparable(e) {
			return e
		}
	}

	return *new(T)
}

// FindList 从数组中返回所有符合条件的项
func FindList[T interface{}](slice []T, comparable func(t T) bool) []T {
	if slice == nil {
		return make([]T, 0)
	}
	var result = make([]T, 0)
	for _, e := range slice {
		if comparable(e) {
			result = append(result, e)
		}
	}

	return result
}

// Select 选择单项，返回集合
func Select[T interface{}, TResult interface{}](slice []T, selectFunc func(t T) TResult) []TResult {
	var result = make([]TResult, 0)

	if slice == nil {
		return result
	}
	for _, e := range slice {
		result = append(result, selectFunc(e))
	}

	return result
}

// SelectString 选择字符串，返回集合
func SelectString[T interface{}](slice []T, selectFunc func(t T) string) []string {
	var result = make([]string, 0)

	if slice == nil {
		return result
	}
	for _, e := range slice {
		result = append(result, selectFunc(e))
	}

	return result
}

// SelectMany 选择集合，返回合并后的集合
func SelectMany[T interface{}, TResult interface{}](slice []T, selectFunc func(t T) []TResult) []TResult {
	var result = make([]TResult, 0)

	if slice == nil {
		return result
	}

	for _, e := range slice {
		result = append(result, selectFunc(e)...)
	}

	return result
}

// GroupBy 分组
func GroupBy[T interface{}, TResult comparable](slice []T, selector func(t T) TResult) map[TResult][]T {
	if slice == nil {
		return make(map[TResult][]T)
	}

	var m = make(map[TResult][]T, 0)
	for _, e := range slice {
		var key = selector(e)
		v, ok := m[key]
		if ok {
			m[key] = append(v, e)
		} else {
			m[key] = []T{e}
		}
	}

	return m
}

// All 所有数据都满足才返回true
func All[T interface{}](slice []T, where func(t T) bool) bool {
	if slice == nil {
		return false
	}
	for _, e := range slice {
		if where(e) == false {
			return false
		}
	}
	return true
}

// Any 只要一个满足就返回true
func Any[T interface{}](slice []T, where func(t T) bool) bool {
	if slice == nil {
		return false
	}

	for _, e := range slice {
		if where(e) {
			return true
		}
	}
	return false
}

func Sum[T interface{}, TTarget interface {
	int | int32 | int64 | float32 | float64
}](slice []T, selectFunc func(t T) TTarget) TTarget {
	if len(slice) == 0 {
		return 0
	}

	var total TTarget = 0
	for _, t := range slice {
		total += selectFunc(t)
	}
	return total
}

// MaxValue 从数组中返回最大值
func MaxValue[T interface{}, TResult interface {
	int | int32 | int64 | float32 | float64
}](slice []T, comparable func(t T) TResult) TResult {
	_, value := comparableMethod(slice, true, comparable)
	return value
}

// MaxOne 从数组中返回最大值的项
func MaxOne[T interface{}, TResult interface {
	int | int32 | int64 | float32 | float64
}](slice []T, comparable func(t T) TResult) T {
	item, _ := comparableMethod(slice, true, comparable)
	return item
}

// MinValue 从数组中返回最小值
func MinValue[T interface{}, TResult interface {
	int | int32 | int64 | float32 | float64
}](slice []T, comparable func(t T) TResult) TResult {
	_, value := comparableMethod(slice, false, comparable)
	return value
}

// MinOne 从数组中返回最小值的项
func MinOne[T interface{}, TResult interface {
	int | int32 | int64 | float32 | float64
}](slice []T, comparable func(t T) TResult) T {
	item, _ := comparableMethod(slice, false, comparable)
	return item
}

// 数值大小比较函数 内部使用
func comparableMethod[T interface{}, TResult interface {
	int | int32 | int64 | float32 | float64
}](slice []T, greater bool, comparable func(t T) TResult) (T, TResult) {
	if len(slice) == 0 {
		return *new(T), 0
	}

	var value TResult
	var item *T
	for i, t := range slice {
		var targetValue = comparable(t)
		if i == 0 {
			value = targetValue
			item = &t
			continue
		}

		if greater {
			if targetValue > value {
				value = targetValue
				item = &t
				continue
			}
		} else {
			if targetValue < value {
				value = targetValue
				item = &t
				continue
			}
		}
	}

	return *item, value
}

// Intersect 集合取交集
func Intersect[T comparable](a []T, b []T) []T {
	var result = make([]T, 0)

	if a == nil || b == nil {
		return result
	}

	for _, e1 := range a {
		if Contains(b, func(item T) bool { return item == e1 }) {
			result = append(result, e1)
		}
	}

	return result
}

// CheckChanged 比较两个集合，返回 新增、修改、删除 的数据
func CheckChanged[T interface{}](source []T, target []T, where func(t1 T, t2 T) bool) ([]T, []T, []T) {
	if source == nil {
		source = make([]T, 0)
	}

	if target == nil {
		target = make([]T, 0)
	}

	var add = make([]T, 0)
	var update = make([]T, 0)
	var deleted = make([]T, 0)
	for _, item := range source {
		if Contains(target, func(t T) bool { return where(t, item) }) {
			update = append(update, item)
		} else {
			add = append(add, item)
		}
	}

	for _, item := range target {
		if Contains(source, func(t T) bool { return where(t, item) }) {

		} else {
			deleted = append(deleted, item)
		}
	}

	return add, update, deleted
}

// Remove 删除集合中指定元素，返回删除后的数据
func Remove[T interface{}](source []*T, where func(item *T) bool) []*T {
	j := 0
	for _, v := range source {
		if where(v) == false {
			source[j] = v
			j++
		}
	}
	return source[:j]
}

func OrderByNumber[T interface{}, TResult interface {
	int | int32 | int64 | float32 | float64
}](slice []T, asc bool, order func(t T) TResult) []T {
	if len(slice) == 0 {
		return make([]T, 0)
	}

	sort.Slice(slice, func(i, j int) bool {
		if asc {
			return order(slice[i]) < order(slice[j])
		} else {
			return order(slice[i]) >= order(slice[j])
		}
	})

	return slice
}

func OrderByString[T interface{}](slice []T, asc bool, order func(t T) string) []T {
	if len(slice) == 0 {
		return make([]T, 0)
	}

	sort.Slice(slice, func(i, j int) bool {
		if asc {
			return order(slice[i]) < order(slice[j])
		} else {
			return order(slice[i]) >= order(slice[j])
		}
	})

	return slice
}

//
//func isPointer(v interface{}) bool {
//	rv := reflect.ValueOf(v)
//	return rv.Kind() == reflect.Ptr
//}
//
//func isPointerArray(arr interface{}) bool {
//	isArray := reflect.TypeOf(arr).Kind() == reflect.Array || reflect.TypeOf(arr).Kind() == reflect.Slice
//	if isArray {
//		return reflect.TypeOf(arr).Elem().Kind() == reflect.Ptr
//	} else {
//		return false
//	}
//}
//
//func isArray(v interface{}) bool {
//	return reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array
//}

func Int64ValuePoint(value int64) *int64 {
	var v = value
	return &v
}

func Int32ValuePoint(value int32) *int32 {
	var v = value
	return &v
}

func IntValuePoint(value int) *int {
	var v = value
	return &v
}

func Float64ValuePoint(value float64) *float64 {
	var v = value
	return &v
}

func Float32ValuePoint(value float32) *float32 {
	var v = value
	return &v
}

func NumberValuePoint[T interface {
	int | int32 | int64 | float32 | float64
}](value T) *T {
	var v = value
	return &v
}

func StringValuePoint(value string) *string {
	if value == "" {
		return nil
	}
	var v = value
	return &v
}

// LocalPager 本地分页
func LocalPager[T interface{}](source []T, page int32, pageSize int32) []T {
	if len(source) == 0 {
		return make([]T, 0)
	}

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	//翻页
	var total = int32(len(source))
	var begin = (page - 1) * pageSize
	var end = pageSize + begin

	if begin >= total {
		return make([]T, 0)
	}

	if end >= total {
		end = total
	}

	return source[begin:end]
}
