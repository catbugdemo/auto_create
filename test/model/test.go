package model

import (
	"fmt"
	"github.com/catbugdemo/errors"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
	"strings"
	"time"
	models "zonst/szghmodels"
)

type SmbVocherDetailLog struct {
	Id                int       `db:"id" json:"id" form:"id"`
	CreateTime        time.Time `db:"create_time" json:"create_time" form:"create_time"`
	CreateDate        time.Time `db:"create_date" json:"create_date" form:"create_date"`
	SmbVocherId       int       `db:"smb_vocher_id" json:"smb_vocher_id" form:"smb_vocher_id"`                      // 卡券种类id
	SmbUserId         int       `db:"smb_user_id" json:"smb_user_id" form:"smb_user_id"`                            // 领取人 用户 id
	SmbVocherDetailId int       `db:"smb_vocher_detail_id" json:"smb_vocher_detail_id" form:"smb_vocher_detail_id"` // 卡券详情id
	VocherDetailCode  string    `db:"vocher_detail_code" json:"vocher_detail_code" form:"vocher_detail_code"`       // 卡券编码
	VocherDetailState int       `db:"vocher_detail_state" json:"vocher_detail_state" form:"vocher_detail_state"`    // 卡券状态 1-未领取 2-已领取 3-已核销 4-已过期
	ReceivedTime      string    `db:"received_time" json:"received_time" form:"received_time"`                      // 领取时间
	VocherType        string    `db:"vocher_type" json:"vocher_type" form:"vocher_type"`                            // 卡券种类
}

func (m *SmbVocherDetailLog) Table() string {
	return "smb_vocher_detail_log"
}

func (m *SmbVocherDetailLog) Condition(arg ...string) string {

	return ""

}

func (m *SmbVocherDetailLog) Insert(db *sqlx.DB) error {
	return Insert(db, m)
}

func (m *SmbVocherDetailLog) Count(db *sqlx.DB, arg ...string) (int, error) {
	return Count(db, m, arg...)
}
func (m *SmbVocherDetailLog) Find(db *sqlx.DB, arg ...string) ([]SmbVocherDetailLog, error) {
	var list []SmbVocherDetailLog
	if err := Find(db, &list, arg...); err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

func (m *SmbVocherDetailLog) First(db *sqlx.DB, arg ...string) error {
	// limit , offset
	arg = append(arg, "1")
	find, err := m.Find(db, arg...)
	if err != nil {
		return errors.WithStack(err)
	}
	*m = find[0]
	return nil
}

func (m *SmbVocherDetailLog) Update(db *sqlx.DB, arg ...string) error {
	return Update(db, m, arg...)
}

func (m *SmbVocherDetailLog) Delete(db *sqlx.DB, arg ...string) error {
	return Delete(db, m, arg...)
}

func (m *SmbVocherDetailLog) IfExist(db *sqlx.DB, arg string) error {
	count, err := m.Count(db, arg)
	if err != nil {
		return errors.WithStack(err)
	}
	if count == 0 {
		return models.ErrRecordNotFound
	}
	return nil
}

func (m *SmbVocherDetailLog) FirstById(db *sqlx.DB, id int) error {
	condition := fmt.Sprintf("where id=%d", id)
	if err := m.IfExist(db, condition); err != nil {
		return errors.WithStack(err)
	}
	return m.First(db, condition)
}

func (m *SmbVocherDetailLog) UpdateById(db *sqlx.DB, id int) error {
	condition := fmt.Sprintf("where id=%d", id)
	if err := m.IfExist(db, condition); err != nil {
		return errors.WithStack(err)
	}

	return m.Update(db, condition)
}

func (m *SmbVocherDetailLog) DeleteById(db *sqlx.DB, id int) error {
	condition := fmt.Sprintf("where id=%d", id)
	if err := m.IfExist(db, condition); err != nil {
		return errors.WithStack(err)
	}
	return m.Delete(db, condition)
}

func (m *SmbVocherDetailLog) FindByCount(db *sqlx.DB,, arg ...string) ([]SmbVocherDetailLog, int, error) {
	var conditon string
	if len(arg) > 0 {
		conditon = arg[0]
	}
	count, err := m.Count(db, conditon)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	if count == 0 {
		return nil, 0, nil
	}
	list, err := m.Find(db, arg...)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return list, count, nil
}

////

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
		params = append(params, fmt.Sprintf("$%d", i))
	}
	return strings.Join(params, ",")
}

func ReturnEqual(params []string) string {
	tmp := make([]string, 0, len(params))
	for i, param := range params {
		tmp = append(tmp, fmt.Sprintf("%s=$%d", param, i+1))
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
