package test

import (
	"encoding/json"
	"fmt"
	"github.com/catbugdemo/auto_create/test/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type SmbWechatPayReturn struct {
	Id            int             `gorm:"column:id;default:" json:"id" form:"id" db:"id"`                                                 //
	CreateTime    time.Time       `gorm:"column:create_time;default:" json:"create_time" form:"create_time" db:"create_time"`             //
	CreateDate    time.Time       `gorm:"column:create_date;default:" json:"create_date" form:"create_date" db:"create_date"`             //
	OutTradeNo    string          `gorm:"column:out_trade_no;default:" json:"out_trade_no" form:"out_trade_no" db:"out_trade_no"`         // 商户系统内部订单号
	TransactionId string          `gorm:"column:transaction_id;default:" json:"transaction_id" form:"transaction_id" db:"transaction_id"` // 微信订单号
	Info          json.RawMessage `gorm:"column:info;default:" json:"info" form:"info" db:"info"`                                         // 存储信息
}

func (m *SmbWechatPayReturn) Table() string {
	return "smb_wechat_pay_return"
}

func (m *SmbWechatPayReturn) Condition(arg ...string) string {

	var params []string
	params = append(params, arg...)

	if m.Id != 0 {
		params = append(params, fmt.Sprint("id='", m.Id, "'"))
	}

	if m.OutTradeNo != "" {
		params = append(params, fmt.Sprint("out_trade_no='", m.OutTradeNo, "'"))
	}

	if m.TransactionId != "" {
		params = append(params, fmt.Sprint("transaction_id='", m.TransactionId, "'"))
	}

	if m.Info != nil {
		params = append(params, fmt.Sprint("info='", m.Info, "'"))
	}

	if len(params) == 0 {
		return ""
	}
	return fmt.Sprintf("where %!s(MISSING)", strings.Join(params, " and "))

}

func (m *SmbWechatPayReturn) Insert(db *sqlx.DB) error {
	return model.Insert(db, m)
}

func (m *SmbWechatPayReturn) Find(db *sqlx.DB, arg ...string) ([]SmbWechatPayReturn, error) {
	var list []SmbWechatPayReturn
	if err := model.Find(db, &list, arg...); err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

func (m *SmbWechatPayReturn) First(db *sqlx.DB, arg ...string) error {
	// limit , offset
	arg = append(arg, "1")
	find, err := m.Find(db, arg...)
	if err != nil {
		return errors.WithStack(err)
	}
	*m = find[0]
	return nil
}

func (m *SmbWechatPayReturn) Update(db *sqlx.DB, arg ...string) error {
	return model.Update(db, m, arg...)
}

func (m *SmbWechatPayReturn) Delete(db *sqlx.DB, arg ...string) error {
	return model.Delete(db, m, arg...)
}

func (m *SmbWechatPayReturn) FirstById(db *sqlx.DB, id int) error {
	return m.First(db, fmt.Sprintf("where id=%!d(MISSING)", id))
}

func (m *SmbWechatPayReturn) UpdateById(db *sqlx.DB, id int) error {
	return m.Update(db, fmt.Sprintf("where id=%!d(MISSING)", id))
}

func (m *SmbWechatPayReturn) DeleteById(db *sqlx.DB, id int) error {
	return m.Delete(db, fmt.Sprintf("where id=%!d(MISSING)", id))
}
