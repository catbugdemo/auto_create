package utils

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var ErrRecordNotFound = errors.New("record not found")

func Delete(db *sqlx.DB, str interface{}, arg ...string) error {
	var where string
	if len(arg) > 0 {
		where = arg[0]
	}

	of := reflect.ValueOf(str)
	if of.Kind() != reflect.Ptr {
		return errors.New("it should be a pointer")
	}
	methodByName := of.MethodByName("Table").Call([]reflect.Value{})

	query := fmt.Sprintf("delete from %s %s",
		methodByName[0].String(),
		where)
	log.Printf("%s", query)
	if _, err := db.Exec(query); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func Update(db *sqlx.DB, str interface{}, arg ...string) error {
	var where string
	if len(arg) > 0 {
		where = arg[0]
	}

	of := reflect.ValueOf(str)
	if of.Kind() != reflect.Ptr {
		return errors.New("it should be a pointer")
	}
	model := reflect.Indirect(of)
	// 转换为结构体，和相应的值
	params, value, err := ReflectInsertDb(model)
	if err != nil {
		return errors.WithStack(err)
	}
	// 插入语句
	methodByName := of.MethodByName("Table").Call([]reflect.Value{})

	query := fmt.Sprintf("update %s set %s %s",
		methodByName[0].String(),
		ReturnEqual(params),
		where,
	)
	log.Printf("%s | %v", query, value)
	if _, err = db.Exec(query, value...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

//
func Find(db *sqlx.DB, str interface{}, arg ...string) error {
	if len(arg) > 3 {
		return errors.New("arg should only where , limit , offset")
	}
	var where, limit, offset string
	if len(arg) > 0 {
		where = arg[0]
	}
	if len(arg) > 1 {
		limit = fmt.Sprintf("limit %s", arg[1])
	}
	if len(arg) > 2 {
		limit = fmt.Sprintf("offset %s", arg[2])
	}

	destSlice := reflect.Indirect(reflect.ValueOf(str))
	destType := destSlice.Type().Elem()
	valueOf := reflect.New(destType)

	// 转换结构体
	params, _, err := ReflectFindDb(reflect.New(destType).Elem())
	if err != nil {
		return errors.WithStack(err)
	}

	query := fmt.Sprintf("select %s from %s %s %s %s",
		strings.Join(params, ","),
		valueOf.MethodByName("Table").Call([]reflect.Value{})[0].String(),
		where,
		limit,
		offset)
	log.Printf("%s", query)
	if err = db.Select(destSlice.Addr().Interface(), query); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func Count(db *sqlx.DB, str interface{}, arg ...string) (int, error) {
	var where string
	if len(arg) > 0 {
		where = arg[0]
	}

	of := reflect.ValueOf(str)
	if of.Kind() != reflect.Ptr {
		return 0, errors.New("it should be a pointer")
	}
	methodByName := of.MethodByName("Table").Call([]reflect.Value{})

	query := fmt.Sprintf("select count(*) from %s %s",
		methodByName[0].String(),
		where)
	log.Printf("%s", query)
	var count int
	if err := db.Get(&count, query); err != nil {
		return 0, errors.WithStack(err)
	}
	return count, nil
}

// 1.判断哪些参数存在
// 2.添加相应参数
func Insert(db *sqlx.DB, str interface{}) error {
	of := reflect.ValueOf(str)
	if of.Kind() != reflect.Ptr {
		return errors.New("it should be a pointer")
	}
	model := reflect.Indirect(of)
	// 转换为结构体，和相应的值
	params, value, err := ReflectInsertDb(model)
	if err != nil {
		return errors.WithStack(err)
	}
	// 插入语句
	methodByName := of.MethodByName("Table").Call([]reflect.Value{})

	query := fmt.Sprintf("insert into %s (%s) values (%s)",
		methodByName[0].String(),
		strings.Join(params, ","),
		ReturnDolor(len(params)))
	log.Printf("%s | %v", query, value)
	if _, err = db.Exec(query, value...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// 返回 db 数据
func ReflectInsertDb(model reflect.Value) ([]string, []interface{}, error) {
	modelType := model.Type()
	params, values := make([]string, 0, model.NumField()), make([]interface{}, 0, model.NumField())
	for i := 0; i < model.NumField(); i++ {
		if IsBlank(model.Field(i)) {
			continue
		}
		// add params
		p := modelType.Field(i)
		if v, ok := p.Tag.Lookup("db"); ok {
			params = append(params, v)
			values = append(values, model.Field(i).Interface())
		}
	}
	return params, values, nil
}

func ReflectFindDb(model reflect.Value) ([]string, []interface{}, error) {
	modelType := model.Type()
	params, values := make([]string, 0, model.NumField()), make([]interface{}, 0, model.NumField())
	for i := 0; i < model.NumField(); i++ {
		// add params
		p := modelType.Field(i)
		if v, ok := p.Tag.Lookup("db"); ok {
			params = append(params, v)
			values = append(values, model.Field(i).Interface())
		}
	}
	return params, values, nil
}

func ReturnDolor(length int) string {
	params := make([]string, 0, length)
	for i := 1; i <= length; i++ {
		// postgres
		params = append(params, fmt.Sprintf("$%d", i))
		// mysql
		// params = append(params, "?")
	}
	return strings.Join(params, ",")
}

func ReturnEqual(params []string) string {
	tmp := make([]string, 0, len(params))
	for i, param := range params {
		// postgres
		tmp = append(tmp, fmt.Sprintf("%s=$%d", param, i+1))
		// mysql
		// tmp = append(tmp, fmt.Sprintf("%s=?", param))
	}
	return strings.Join(tmp, ",")
}

func IsBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// 获取 limit, offset
func ToLimitOffset(sizeIn string, indexIn string, count int) (limit int, offset int) {
	size, _ := strconv.Atoi(sizeIn)
	index, _ := strconv.Atoi(indexIn)
	if size == 0 {
		size = 10
	}

	if index == 0 {
		index = 1
	}
	//1
	if count == 0 {
		return size, 0
	}
	var pageMax int
	//1%10
	if count%size == 0 {
		pageMax = count / size
	} else {
		//1
		pageMax = count/size + 1
	}
	//1<=1
	if pageMax <= index {
		index = pageMax
	}
	offset = size * (index - 1)

	if offset == -10 {
		offset = 0
	}
	return size, offset
}
