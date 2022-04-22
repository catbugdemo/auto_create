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
		DataSource: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "1.117.233.151", "5432", "ytf", "smb", "disable", "ytf@2021"),
		TableName:  "smb_wechat_pay_return",
		Driver:     "postgres",
	}

	generate, err := AutoGenerateModel(&normal)
	assert.Nil(t, err)
	fmt.Printf(generate)
}

func TestSqlx(t *testing.T) {
	normal := Normal{
		DataSource: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "1.117.233.151", "5432", "ytf", "smb", "disable", "ytf@2021"),
		TableName:  "smb_wechat_pay_return",
		Driver:     "postgres",
	}

	generate, err := AutoGenerateSqlx(&normal)
	assert.Nil(t, err)
	fmt.Printf(generate)
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
		ModelsName:  "model.YtfAdvAccountInfo",
		RedisConfig: `c.MustGet(REDIS_TOKEN).(*redis.Pool).Get()`,
		Handlers:    "handlers",
	}

	fmt.Printf(AutoGenerateCRUD(&st))
}

func TestSqlxCrud(t *testing.T) {
	type Str struct {
	}
	st := St{
		Stru:     Str{},
		DbConfig: `c.MustGet(DB_CONFIG).(*gorm.DB)`,
		Info: map[string]interface{}{
			"bind_json":    "",
			"service_json": "",
			"success_json": "",
		},
	}
	fmt.Print(AutoGenerateSqlxControl(&st))
}

func TestModelConvert(t *testing.T) {
	var sql = `
create table smb_province_code(
  id serial primary key not null, 
  create_time timestamp with time zone DEFAULT now(),
	create_date date NULL DEFAULT 'now'::text::date,
  
  province varchar not null default '', -- 省
  city varchar not null default '', -- 市
  area varchar not null default ''-- 区
  code int not null default 0,  -- 编码
);
create index on smb_province_code(province,city,area);
`

	fmt.Println(model_convert.GenerateNote(sql))
}
