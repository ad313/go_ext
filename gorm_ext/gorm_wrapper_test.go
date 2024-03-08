package gorm_ext

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type Table3 struct {
	Id           string         `gorm:"column:id;type:VARCHAR2(36);primaryKey;not null"` //标识
	Name         string         `gorm:"column:name;type:VARCHAR2(20)" json:"name"`
	Sn           string         `gorm:"column:sn;type:VARCHAR2(50)" json:"sn"`
	LinkSn       string         `gorm:"column:link_sn;type:VARCHAR2(50)" json:"link_sn"`
	Mode         string         `gorm:"column:mode;type:VARCHAR2(50)" json:"mode"`
	Domain       int32          `gorm:"column:domain;type:INTEGER(4)" json:"domain"`
	Type         int32          `gorm:"column:type;type:INTEGER(4)" json:"type"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME(8)" json:"deleted_at"`
	CreationTime int64          `gorm:"column:creation_time;type:BIGINT(8)" json:"creation_time"`
	CreatorId    string         `gorm:"column:creator_id;type:VARCHAR2(36)" json:"creator_id"`
}

type Table4 struct {
	Id           string         `gorm:"column:id;type:VARCHAR2(36);primaryKey;not null"` //标识
	Name         string         `gorm:"column:name;type:VARCHAR2(20)" json:"name"`
	Sn           string         `gorm:"column:sn;type:VARCHAR2(50)" json:"sn"`
	LinkSn       string         `gorm:"column:link_sn;type:VARCHAR2(50)" json:"link_sn"`
	Mode         string         `gorm:"column:mode;type:VARCHAR2(50)" json:"mode"`
	Domain       int32          `gorm:"column:domain;type:INTEGER(4)" json:"domain"`
	Type         int32          `gorm:"column:type;type:INTEGER(4)" json:"type"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME(8)" json:"deleted_at"`
	CreationTime int64          `gorm:"column:creation_time;type:BIGINT(8)" json:"creation_time"`
	CreatorId    string         `gorm:"column:creator_id;type:VARCHAR2(36)" json:"creator_id"`
}

type Table5 struct {
	Id           string         `gorm:"column:id;type:VARCHAR2(36);primaryKey;not null"` //标识
	Name         string         `gorm:"column:name;type:VARCHAR2(20)" json:"name"`
	Sn           string         `gorm:"column:sn;type:VARCHAR2(50)" json:"sn"`
	LinkSn       string         `gorm:"column:link_sn;type:VARCHAR2(50)" json:"link_sn"`
	Mode         string         `gorm:"column:mode;type:VARCHAR2(50)" json:"mode"`
	Domain       int32          `gorm:"column:domain;type:INTEGER(4)" json:"domain"`
	Type         int32          `gorm:"column:type;type:INTEGER(4)" json:"type"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME(8)" json:"deleted_at"`
	CreationTime int64          `gorm:"column:creation_time;type:BIGINT(8)" json:"creation_time"`
	CreatorId    string         `gorm:"column:creator_id;type:VARCHAR2(36)" json:"creator_id"`
}

func (t *Table3) TableName() string {
	return "Table3"
}
func (t *Table4) TableName() string {
	return "Table4"
}
func (t *Table5) TableName() string {
	return "Table5"
}

var mysqlConn = "root:123456@tcp(192.168.1.80:30680)/test?charset=utf8mb4&parseTime=True&loc=Local"

func init() {
	Init(NewDb(mysqlConn), MySql)
}

var table3 = BuildGormTable[Table3]().Table.T
var table4 = BuildGormTable[Table4]().Table.T
var table5 = BuildGormTable[Table5]().Table.T

// GetDbContext 获取DbContext。当外部开启事务时，传入开启事务后的db
func (a *Table3) GetDbContext(ctx context.Context, db ...*gorm.DB) *OrmWrapper[Table3] {
	return BuildOrmWrapper[Table3](ctx, db...)
}

func NewDb(conn string) *gorm.DB {
	if len(conn) == 0 {
		panic("数据库连接字符未配置")
	}

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic("创建mysql 数据库失败")
	}
	return db
}

func Test_Wrapper(t *testing.T) {

	sql, err := table3.GetDbContext(context.Background()).
		SetTableAlias("a").
		//SelectColumnOriginal("a.*", "").
		//SelectColumnOriginal("count(*)", "aaa").
		WhereCondition(&Condition{
			Column:         &table3.Id,
			CompareSymbols: Eq,
			Arg:            "123",
			TableAlias:     "a",
		}).
		WhereByColumn(&table3.Name, Like, "abc", "a").
		WhereByColumn(&table3.Name, StartWith, "abc", "b").
		WhereByColumn(&table3.Name, EndWith, "abc", "c").
		WhereCondition(&ExistsCondition{
			Table: table3,
			ConditionBuilder: NewAndConditionBuilder(&TableCondition{
				InnerColumn:    &table3.Id,
				OuterAlias:     "a",
				OuterColumn:    &table3.Id,
				CompareSymbols: Eq,
			}, &TableCondition{
				InnerColumn:    &table3.Name,
				OuterAlias:     "a",
				OuterColumn:    &table3.Name,
				CompareSymbols: NotEq,
			}, &Condition{
				TableAlias:     "",
				Column:         &table3.Type,
				CompareSymbols: Gt,
				Arg:            18,
			}),
			IsNotExists: true,
			Func:        "",
		}).
		LeftJoin(table4, "b", &table3.Id, &table4.Id).
		LeftJoin(table5, "c", &table4.Sn, &table5.LinkSn).
		//Unscoped().
		ToSql()

	if err != nil {
		t.Errorf("Test_Wrapper faild")
	}

	fmt.Println(sql)

	table3.GetDbContext(context.Background()).
		LeftJoin(table4, "b", &table3.Id, &table4.Id).
		ToPagerList(&Pager{Page: 0, PageSize: 0})

}
