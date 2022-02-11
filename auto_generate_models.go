package auto

import (
	"bytes"
	"fmt"
	"html/template"
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
	Package    string                 `json:"package"`
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
	if n.Package == "" {
		n.Package = "models"
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

	/*	var str = `

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
	*/
	var str2 = `
package {{.package}}

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type {{.struct_name}} struct {
	${type_struct}
}

// RedisKey
func (o *{{.struct_name}}) RedisKey() string {
	// TODO set redis key
	return fmt.Sprintf()
}

func (o *{{.struct_name}}) ArrayRedisKey() string {
	// TODO set array redis key
	return fmt.Sprintf()
}

// RedisDuration
func (o *{{.struct_name}}) RedisDuration() time.Duration {
	// TODO set redis duration , default 30 ~ 60 minutes
	return time.Duration((rand.Intn(60-30) + 30)) * time.Minute
}

// SyncToRedis
func (o *{{.struct_name}}) SyncToRedis(conn *redis.Conn) error {
	if o.RedisKey() == "" {
		return errors.New("not set redis key")
	}
	buf, err := json.Marshal(o)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = conn.SetEX(context.Background(), o.RedisKey(), string(buf), o.RedisDuration()).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetFromRedis
func (o *{{.struct_name}}) GetFromRedis(conn *redis.Conn) error {
	if o.RedisKey() == "" {
		return errors.New("not set redis key")
	}
	buf, err := conn.Get(context.Background(), o.RedisKey()).Bytes()
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		return errors.WithStack(err)
	}
	if string(buf) == "DISABLE" {
		return errors.New("not found data in redis nor db")
	}

	if err = json.Unmarshal(buf, o); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (o *{{.struct_name}}) DeleteFromRedis(conn *redis.Conn) error {
	if o.RedisKey() != "" {
		if err := conn.Del(context.Background(), o.RedisKey()).Err(); err != nil {
			return errors.WithStack(err)
		}
	}
	if o.ArrayRedisKey() != "" {
		if err := conn.Del(context.Background(), o.ArrayRedisKey()).Err(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// MustGet
func (o *{{.struct_name}}) MustGet(engine *gorm.DB, conn *redis.Conn) error {
	err := o.GetFromRedis(conn)
	if err == nil {
		return nil
	}

	if err != nil && err != redis.Nil {
		return errors.WithStack(err)
	}
	// not found in redis
	var count int64
	if err = engine.Count(&count).Error; err != nil {
		return errors.WithStack(err)
	}
	// prevent cache penetration
	if count == 0 {
		if err = conn.SetNX(context.Background(), o.RedisKey(), "DISABLE", o.RedisDuration()).Err(); err != nil {
			return errors.WithStack(err)
		}
		return errors.New("not found data in redis nor db")
	}

	var mutex = o.RedisKey() + "_MUTEX"
	if err = conn.Get(context.Background(), mutex).Err(); err != nil {
		if err != redis.Nil {
			return errors.WithStack(err)
		}
		// set redis mutex and get from db
		if err = conn.SetNX(context.Background(), mutex, 1, 5*time.Second).Err(); err != nil {
			return errors.WithStack(err)
		}
		if err = engine.First(&o).Error; err != nil {
			return errors.WithStack(err)
		}
		if err = o.SyncToRedis(conn); err != nil {
			return errors.WithStack(err)
		}
		if err = conn.Del(context.Background(), mutex).Err(); err != nil {
			return errors.WithStack(err)
		}
	} else {
		// found lock , waiting unlock
		var index int
		for {
			if index > 10 {
				return errors.New(mutex + " lock error")
			}
			if err2 := conn.Get(context.Background(), mutex).Err(); err2 != nil {
				break
			} else {
				time.Sleep(30 * time.Millisecond)
				index++
				continue
			}
		}
		if err = o.MustGet(engine, conn); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// ArraySyncToRedis
func (o *{{.struct_name}}) ArraySyncToRedis(list []{{.struct_name}}, conn *redis.Conn) error {
	if o.ArrayRedisKey() == "" {
		return errors.New("not set redis key")
	}
	buf, err := json.Marshal(list)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = conn.SetEX(context.Background(), o.ArrayRedisKey(), string(buf), o.RedisDuration()).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ArrayGetFromRedis
func (o *{{.struct_name}}) ArrayGetFromRedis(conn *redis.Conn) ([]{{.struct_name}}, error) {
	if o.RedisKey() == "" {
		return nil, errors.New("not set redis key")
	}
	buf, err := conn.Get(context.Background(), o.ArrayRedisKey()).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, redis.Nil
		}
		return nil, errors.WithStack(err)
	}
	if string(buf) == "DISABLE" {
		return nil, errors.New("not found data in redis nor db")
	}
	var list []{{.struct_name}}
	if err = json.Unmarshal(buf, &list); err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

// ArrayDeleteFromRedis
func (o *{{.struct_name}}) ArrayDeleteFromRedis(conn *redis.Conn) error {
	return o.DeleteFromRedis(conn)
}

// ArrayMustGet
func (o *{{.struct_name}}) ArrayMustGet(engine *gorm.DB, conn *redis.Conn) ([]{{.struct_name}}, error) {
	list, err := o.ArrayGetFromRedis(conn)
	if err == nil {
		return list, nil
	}
	if err != nil && err != redis.Nil {
		return nil, errors.WithStack(err)
	}

	// not found in redis
	var count int64
	if err = engine.Count(&count).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	// prevent cache penetration
	if count == 0 {
		if err = conn.SetNX(context.Background(), o.ArrayRedisKey(), "DISABLE", o.RedisDuration()).Err(); err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.New("not found data in redis nor db")
	}

	var mutex = o.ArrayRedisKey() + "_MUTEX"
	if err = conn.Get(context.Background(), mutex).Err(); err != nil {
		if err != redis.Nil {
			return nil, errors.WithStack(err)
		}
		// set redis mutex and get from db
		if err = conn.SetNX(context.Background(), mutex, 1, 5*time.Second).Err(); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = engine.Find(&list).Error; err != nil {
			return nil, errors.WithStack(err)
		}
		if err = o.ArraySyncToRedis(list, conn); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = conn.Del(context.Background(), mutex).Err(); err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		// found lock , waiting unlock
		var index int
		for {
			if index > 10 {
				return nil, errors.New(mutex + " lock error")
			}
			if err2 := conn.Get(context.Background(), mutex).Err(); err2 != nil {
				break
			} else {
				time.Sleep(30 * time.Millisecond)
				index++
				continue
			}
		}
		list, err = o.ArrayMustGet(engine, conn)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return list, nil
}
`
	str2 = strings.ReplaceAll(str2, "${type_struct}", getTypeStruct(columns))
	tt := template.Must(template.New("model").Parse(str2))
	vals := map[string]string{
		"package":     n.Package,
		"struct_name": n.Info["struct_name"].(string),
		//"type_struct": getTypeStruct(columns),
	}
	var buf bytes.Buffer
	if err := tt.Execute(&buf, vals); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getTypeStruct(columns []Column) string {
	var tmp string
	for _, column := range columns {
		tmp += fmt.Sprintf("    %s  %s    `gorm:\"column:%s;default:\" json:\"%s\" form:\"%s\" db:\"%s\"` // %s \n",
			underLineToHump(column.ColumnName), typeConvert(column.ColumnType), column.ColumnName, column.ColumnName, column.ColumnName, column.ColumnName, column.ColumnComment)
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
	ColumnName    string `gorm:"column:column_name"`    // column_name
	ColumnType    string `gorm:"column:column_type"`    // column_type
	ColumnComment string `gorm:"column:column_comment"` // column_comment
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
			col_description(a.attrelid, a.attnum) as column_comment,
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
