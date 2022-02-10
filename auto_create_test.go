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
		TableName:  "smb_tencentad_err_log",
		Driver:     "postgres",
	}

	generate, err := AutoGenerateModel(&normal)
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
		ModelsName:  "models.YtfAdvAccountInfo",
		RedisConfig: `c.MustGet(REDIS_TOKEN).(*redis.Pool).Get()`,
		Handlers:    "handlers",
	}

	fmt.Printf(AutoGenerateCRUD(&st))
}

func TestModelConvert(t *testing.T) {
	var sql = `
create table smb_tencentad_local_log(
  id serial primary key,
  create_time timestamp with time zone DEFAULT now(),
	create_date date DEFAULT 'now'::text::date,
  
  access_token varchar not null default '', -- access_token
  account_id int not null default 0, -- 腾讯 account_id 
  req jsonb, -- 请求参数
  resp jsonb -- 返回值
);
create index on smb_tencentad_local_log(account_id);
`

	fmt.Println(model_convert.GenerateNote(sql))
}
