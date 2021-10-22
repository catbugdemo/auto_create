package auto

import (
	"fmt"
	"strings"

	"github.com/catbugdemo/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Normal struct {
	DataSource string                 `json:"data_source"`
	TableName  string                 `json:"table_name"`
	Driver     string                 `json:"driver"`
	Info       map[string]interface{} `json:"info"`
}

type Way interface {
	init() error
	formatJSON() (string, error)
}

func AutoGenerateModel(way Way) (string, error) {
	if err := way.init(); err != nil {
		return "", errors.WithStack(err)
	}

	json, err := way.formatJSON()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return json, nil
}

func (n *Normal) init() error {
	if n.DataSource == "" || n.TableName == "" {
		return errors.New("DataSource or TableName is nill")
	}
	n.initInfo()
	return nil
}

func (n *Normal) initInfo() {
	if n.Info == nil {
		n.Info = map[string]interface{}{}

	}
	n.Info["struct_name"] = underLineToHump(n.TableName)
}

// formatJSON autoCreate
// ${type_struct}  结构体
// ${struct_name} 生成结构体名称
// ${table_name} 表名
func (n *Normal) formatJSON() (string, error) {
	columns := FindColumns(n.Driver, n.DataSource, n.TableName)

	var str = `

// without 2-cache
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"github.com/catbugdemo/errors"

	"github.com/gomodule/redigo/redis"
)

/*
	need pkg
		- "gorm.io/gorm"
		- "github.com/catbugdemo/errors" or "github.com/pkg/errors"
		- "github.com/gomodule/redigo/redis"
*/

type ${struct_name} struct {
	${type_struct}
}

func (o ${struct_name}) TableName() string {
	return "${table_name}"
}

var ${struct_name}RedisKeyFormat = ""

func (o ${struct_name}) RedisKey() string {
	// TODO set redis key account to index
	return fmt.Sprintf(${struct_name}RedisKeyFormat)
}

var Array${struct_name}RedisKeyFormat = ""

func (o ${struct_name}) ArrayRedisKey() string {
	// TODO set its redis key account to index
	return fmt.Sprintf(Array${struct_name}RedisKeyFormat)
}

func (o ${struct_name}) RedisSecondDuration() int {
	// TODO set redis duration default 1 - 7 days , return -1 means not time limit
	return (rand.Intn(7-1) + 1) * 60 * 60 * 24
}

func (o *${struct_name}) GetFromRedis(conn redis.Conn) error {
	if o.RedisKey() == "" {
		return errors.New("not set redis key")
	}

	buf, err := redis.Bytes(conn.Do("GET", o.RedisKey()))

	if err != nil {
		if err == redis.ErrNil {
			return redis.ErrNil
		}
		return errors.WithStack(err)
	}

	// Prevent cache penetration
	if string(buf) == "DISABLE" {
		return errors.New("not found data in redis nor db")
	}

	if err = json.Unmarshal(buf, &o); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o ${struct_name}) ArrayGetFromRedis(conn redis.Conn) ([]${struct_name}, error) {
	if o.ArrayRedisKey() == "" {
		return nil, errors.New("not set redis key")
	}

	buf, err := redis.Bytes(conn.Do("GET", o.RedisKey()))

	if err != nil {
		if err == redis.ErrNil {
			return nil, redis.ErrNil
		}
		return nil, errors.WithStack(err)
	}

	// Prevent cache penetration
	if string(buf) == "DISABLE" {
		return nil, fmt.Errorf("not found data in redis nor db")
	}

	list := make([]${struct_name}, 0)
	if err = json.Unmarshal(buf, &list); err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

// SyncToRedis sync redis
func (o ${struct_name}) SyncToRedis(conn redis.Conn) error {
	if o.RedisKey() == "" {
		return errors.New("not set redis key")
	}

	buf, err := json.Marshal(o)
	if err != nil {
		return errors.WithStack(err)
	}

	if o.RedisSecondDuration() == -1 {
		if _, err := conn.Do("SET", o.RedisKey(), buf); err != nil {
			return errors.WithStack(err)
		}
	} else {
		if _, err := conn.Do("SETEX", o.RedisKey(), o.RedisSecondDuration(), buf); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// MustGet query single data
func (o *${struct_name}) MustGet(db *gorm.DB, conn redis.Conn) error {

	err := o.GetFromRedis(conn)

	if err != nil && err.Error() == "not found data in redis nor db" {
		return errors.WithStack(err)
	}

	if err == redis.ErrNil {
		// get from db
		var count int64
		if err2 := db.Count(&count).Error; err2 != nil {
			return errors.WithStack(err2)
		}
		// Prevent cache penetration
		if count == 0 {
			if o.RedisSecondDuration() == -1 {
				_, _ = conn.Do("SETNX", o.RedisKey(), "DISABLE")
			} else {
				_, _ = conn.Do("SET", o.RedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
			return errors.New("not found data in redis nor db")
		}

		if err3 := db.First(&o).Error; err3 != nil {
			return errors.WithStack(err3)
		}

		if err4 := o.SyncToRedis(conn); err != nil {
			return errors.WithStack(err4)
		}
		return nil
	}

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ArraySyncToRedis
func (o ${struct_name}) ArraySyncToRedis(list []${struct_name}, conn redis.Conn) error {
	if o.ArrayRedisKey() == "" {
		return errors.New("not set redis key")
	}

	buf, err := json.Marshal(list)
	if err != nil {
		return errors.WithStack(err)
	}

	if o.RedisSecondDuration() == -1 {
		if _, err := conn.Do("SET", o.RedisKey(), buf); err != nil {
			return errors.WithStack(err)
		}
	} else {
		if _, err := conn.Do("SETEX", o.RedisKey(), o.RedisSecondDuration(), buf); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// ArrayMustGet query list data
func (o ${struct_name}) ArrayMustGet(db *gorm.DB, conn redis.Conn) ([]${struct_name}, error) {


	list, err := o.ArrayGetFromRedis(conn)
	if err != nil && err.Error() == "not found data in redis nor db" {
		// redis value is DISABLE
		return nil, errors.WithStack(err)
	}

	if err == redis.ErrNil {
		// get from db
		var count int64
		if err2 := db.Count(&count).Error; err2 != nil {
			return nil, errors.WithStack(err2)
		}

		// Prevent cache penetration
		if count == 0 {
			if o.RedisSecondDuration() == -1 {
				_, _ = conn.Do("SETNX", o.RedisKey(), "DISABLE")
			} else {
				_, _ = conn.Do("SET", o.RedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
			return nil, errors.New("not found data in redis nor db")
		}

		if err3 := db.Find(&list).Error; err3 != nil {
			return nil, errors.WithStack(err3)
		}

		if err4 := o.ArraySyncToRedis(list, conn); err4 != nil {
			return nil, errors.WithStack(err)
		}

		return list, nil
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

// DeleteFromRedis delete redis
func (o ${struct_name}) DeleteFromRedis(conn redis.Conn) error {
	if o.RedisKey() != "" {
		if _, err := conn.Do("DEL", o.RedisKey()); err != nil {
			return errors.WithStack(err)
		}
	}

	if o.ArrayRedisKey() != "" {
		if _, err := conn.Do("DEL", o.RedisKey()); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// ArrayDeleteFromRedis  delete list redis
func (o ${struct_name}) ArrayDeleteFromRedis(conn redis.Conn) error {
	return o.DeleteFromRedis(conn)
}


`
	str = strings.ReplaceAll(str, "${type_struct}", getTypeStruct(columns))
	str = strings.ReplaceAll(str, "${struct_name}", n.Info["struct_name"].(string))
	str = strings.ReplaceAll(str, "${table_name}", n.TableName)
	return str, nil
}

func getTypeStruct(columns []Column) string {
	var tmp string
	for _, column := range columns {
		tmp += fmt.Sprintf("    %s  %s    `gorm:\"column:%s;default:\" json:\"%s\" form:\"%s\"`\n",
			underLineToHump(column.ColumnName), typeConvert(column.ColumnType), column.ColumnName, column.ColumnName, column.ColumnName)
	}
	return tmp[:len(tmp)-1]
}

func typeConvert(s string) string {
	if strings.Contains(s, "char") || in(s, []string{
		"text",
	}) {
		return "string"
	}
	// postgres
	{
		if in(s, []string{"double precision", "double"}) {
			return "float64"
		}
		if in(s, []string{"bigint", "bigserial", "integer", "smallint", "serial", "big serial"}) {
			return "int"
		}
		if in(s, []string{"numeric", "decimal", "real"}) {
			return "decimal.Decimal"
		}
		if in(s, []string{"bytea"}) {
			return "[]byte"
		}
		if strings.Contains(s, "time") || in(s, []string{"date", "datetime", "timestamp"}) {
			return "time.Time"
		}
		if in(s, []string{"jsonb"}) {
			return "json.RawMessage"
		}
		if in(s, []string{"bool", "boolean"}) {
			return "bool"
		}

		if in(s, []string{"bigint[]"}) {
			return "[]int64"
		}
	}
	// mysql
	{
		if strings.HasPrefix(s, "int") {
			return "int"
		}
		if strings.HasPrefix(s, "varchar") {
			return "string"
		}
		if s == "json" {
			return "json.RawMessage"
		}
		if in(s, []string{"bool", "boolean"}) {
			return "bool"
		}
	}

	return s
}

// s 是否in arr
func in(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

// UnderLineToHump 下划线转驼峰
func underLineToHump(s string) string {
	arr := strings.Split(s, "_")
	for i, v := range arr {
		arr[i] = strings.ToUpper(string(v[0])) + v[1:]
	}
	return strings.Join(arr, "")
}

// 数据库列属性
type Column struct {
	ColumnName string `gorm:"column:column_name"` // column_name
	ColumnType string `gorm:"column:column_type"` // column_type
}

// 根据数据源，表明获取列属性
func FindColumns(dialect string, dataSource string, tableName string) []Column {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(fmt.Sprintf("recover from a fatal error : %v", e))
		}
	}()

	switch dialect {
	case "postgres", "pg", "psql":
		return findPGColumns(dataSource, tableName)
	case "mysql":
		return findMysqlColumns(dataSource, tableName)
	}
	return nil
}

func findPGColumns(dataSource string, tableName string) []Column {
	var FindColumnsSql = `
        SELECT
            a.attnum AS column_number,
            a.attname AS column_name,
            --format_type(a.atttypid, a.atttypmod) AS column_type,
            a.attnotnull AS not_null,
			COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,
    		COALESCE(ct.contype = 'p', false) AS  is_primary_key,
    		CASE
        	WHEN a.atttypid = ANY ('{int,int8,int2}'::regtype[])
          		AND EXISTS (
				SELECT 1 FROM pg_attrdef ad
             	WHERE  ad.adrelid = a.attrelid
             	AND    ad.adnum   = a.attnum
             	-- AND    ad.adsrc = 'nextval('''
                --	|| (pg_get_serial_sequence (a.attrelid::regclass::text
                --	                          , a.attname))::regclass
                --	|| '''::regclass)'
             	)
            THEN CASE a.atttypid
                    WHEN 'int'::regtype  THEN 'serial'
                    WHEN 'int8'::regtype THEN 'bigserial'
                    WHEN 'int2'::regtype THEN 'smallserial'
                 END
			WHEN a.atttypid = ANY ('{uuid}'::regtype[]) AND COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') != ''
            THEN 'autogenuuid'
        	ELSE format_type(a.atttypid, a.atttypmod)
    		END AS column_type
		FROM pg_attribute a
		JOIN ONLY pg_class c ON c.oid = a.attrelid
		JOIN ONLY pg_namespace n ON n.oid = c.relnamespace
		LEFT JOIN pg_constraint ct ON ct.conrelid = c.oid
		AND a.attnum = ANY(ct.conkey) AND ct.contype = 'p'
		LEFT JOIN pg_attrdef ad ON ad.adrelid = c.oid AND ad.adnum = a.attnum
		WHERE a.attisdropped = false
		AND n.nspname = 'public'
		AND c.relname = ?
		AND a.attnum > 0
		ORDER BY a.attnum
	`
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		panic(err)
	}
	var columns = make([]Column, 0, 10)
	db.Raw(FindColumnsSql, tableName).Find(&columns)
	return columns
}
func findMysqlColumns(dataSource string, tableName string) []Column {
	var FindColumnsSql = `
       SELECT column_name as column_name, column_type as column_type  FROM information_schema.columns WHERE table_name= ?
	`
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var columns = make([]Column, 0, 10)
	if e := db.Raw(FindColumnsSql, tableName).Scan(&columns).Error; e != nil {
		panic(err)
	}

	return columns
}
