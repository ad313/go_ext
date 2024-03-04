package gorm_ext

// 定义数据库类型
const (
	MySql     = "mysql"
	Postgres  = "postgres"
	Sqlite    = "sqlite"
	Sqlserver = "sqlserver"
	Dm        = "dm"
)

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
