package auto

import (
	"fmt"
	"github.com/fwhezfwhez/model_convert"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 生成 Controller 层,包括 swag
func TestGenerateController(t *testing.T) {
	m := Model{
		Control: Control{
			ControlName: "GetAdgroups", // 输入创建名称
			DbConfig:    "c.MustGet(DB_CONFIG).(*gorm.DB)",

			// log
			LogBind:    "JsonNotifyStatusWithLog(c, CODE_FAIL_I, err.Error(), req, timeNow)",
			LogService: "JsonNotifyStatusWithLog(c, CODE_FAIL_I, err.Error(), req, timeNow)",
			LogReturn:  "JsonNotifyRetWithLog(c, CODE_SUCCESS_I, CODE_SUCCESS_S, req, data, timeNow)",
		},

		Swag: Swag{
			Security: "x-smb-jwt",
		},
	}
	fmt.Println(GenerateController(m))
}

// 生成底层模板
func TestModels(t *testing.T) {
	normal := Normal{
		DataSource: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "1.117.233.151", "5432", "ytf", "ytf", "disable", "ytf@2021"),
		TableName:  "ytf_adviser_building",
		Driver:     "postgres",
	}

	generate, err := AutoGenerateModel(&normal)
	assert.Nil(t, err)
	fmt.Printf(generate)
}

func TestSqlx(t *testing.T) {
	normal := Normal{
		//DataSource: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "1.117.233.151", "5432", "ytf", "smb", "disable", "ytf@2021"),
		DataSource: fmt.Sprint("new_retailers:mnHL63mzzNX2GyYd@tcp(1.15.221.224:3306)/new_retailers"),
		TableName:  "adv_shop_online_customer",
		Driver:     "mysql",
	}

	generate, err := AutoGenerateSqlx(&normal)
	assert.Nil(t, err)
	fmt.Println(generate)
}

// 自动生成 crud
func TestCRUD(t *testing.T) {
	type YtfAdvAccountInfo struct {
		Id        int       `gorm:"column:id;default:" json:"id" form:"id"`
		CreatedAt time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`

		Qq              string `gorm:"column:qq;default:" json:"qq" form:"qq"`
		AccountId       int    `gorm:"column:account_id;default:" json:"account_id" form:"account_id"`
		CorporationName string `gorm:"column:corporation_name;default:" json:"corporation_name" form:"corporation_name"`
	}

	st := St{
		Stru:        YtfAdvAccountInfo{},
		DbConfig:    `c.MustGet(DB_CONFIG).(*gorm.DB)`,
		ModelsName:  "models.YtfAdvAccountInfo",
		RedisConfig: `c.MustGet(REDIS_TOKEN).(*redis.Pool).Get()`,
		Handlers:    "handlers",
		Info: map[string]interface{}{
			"tag": "测试",
		},
	}

	fmt.Printf(AutoGenerateCRUD(&st))
}

func TestSqlxCrud(t *testing.T) {
	type AdvShopOnlineCustomer struct {
		Id          int       `gorm:"column:id;default:" json:"id" form:"id" db:"id"`                                             // 自增id
		CreateTime  time.Time `gorm:"column:create_time;default:" json:"create_time" form:"create_time" db:"create_time"`         // 创建时间
		Name        string    `gorm:"column:name;default:" json:"name" form:"name" db:"name"`                                     // 商户名称
		NickName    string    `gorm:"column:nick_name;default:" json:"nick_name" form:"nick_name" db:"nick_name"`                 // 商户昵称
		IconUrl     string    `gorm:"column:icon_url;default:" json:"icon_url" form:"icon_url" db:"icon_url"`                     // 头像url
		IconMediaId string    `gorm:"column:icon_media_id;default:" json:"icon_media_id" form:"icon_media_id" db:"icon_media_id"` // 微信头像id
		ShopId      int       `gorm:"column:shop_id;default:" json:"shop_id" form:"shop_id" db:"shop_id"`                         // 商户id
		BusinessId  int       `gorm:"column:business_id;default:" json:"business_id" form:"business_id" db:"business_id"`         // 商户在线id
		Status      int       `gorm:"column:status;default:" json:"status" form:"status" db:"status"`                             // 1:正常 2:删除
	}
	st := St{
		Stru:     AdvShopOnlineCustomer{},
		DbConfig: `c.MustGet(DB_CONFIG).(*sqlx.DB)`,
		Info: map[string]interface{}{
			"tag": "测试",
		},
	}
	fmt.Println(AutoGenerateSqlxControl(&st))
}

func TestModelConvert(t *testing.T) {
	var sql = `
create table ytf_adviser_building (
  id int not null default 0,
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  
  adviser_id int not null default 0, -- 置业顾问id 
  merchant_id int not null default 0 , -- 商户 id 
  building_id int not null default 0, -- 楼盘id 
)
create unique index on ytf_adviser_building(adviser_id,merchant_id,building_id)
`

	fmt.Println(model_convert.GenerateNote(sql))
}
