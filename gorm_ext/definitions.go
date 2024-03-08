package gorm_ext

// 定义数据库类型
const (
	MySql     = "mysql"
	Postgres  = "postgres"
	Sqlite    = "sqlite"
	Sqlserver = "sqlserver"
	Dm        = "dm"
)

// 定义运算符
const (
	Eq        = "="
	NotEq     = "<>"
	Gt        = ">"
	GtAndEq   = ">="
	Less      = "<"
	LessAndEq = "<="
	In        = "IN"
	NotIn     = "NOT IN"
	Like      = "LIKE"
	NotLike   = "NOT LIKE"
	StartWith = "STARTWITH"
	EndWith   = "ENDWITH"
	IsNull    = "IS NULL"
	NotNull   = "IS NOT NULL"
)

// 聚合函数
const (
	Max   = "Max"   //最大值
	Min   = "Min"   //最小值
	Avg   = "Avg"   //平均值
	Sum   = "Sum"   //求和
	Count = "Count" //统计行
)

// 数值型函数
const (
	Abs   = "Abs"   //求绝对值
	Sqrt  = "Sqrt"  //开平方根
	Ceil  = "Ceil"  //向上取整
	Floor = "Floor" //向下取整
	Round = "Round" //四舍五入
)

// 字符串函数
const (
	Upper = "Upper" //转大写
	Lower = "Lower" //转小写
)
