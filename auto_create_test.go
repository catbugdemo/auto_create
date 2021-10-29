package auto

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 生成 Controller 层
func TestGenerateController(t *testing.T) {
	c := Control{
		ControlName: "GetArrayPeopleTarget", // 输入创建名称
		Describe:    "通过 1 2 级行业搜索定向包",      // 输入描述 -- 可不填

		ReqBool:  true, // 是否需要手动填写绑定参数，推荐3个以内为true
		Req:      "",   // 请求名称,如果 ReqBool == true 不填
		DbConfig: "c.MustGet(DB_CONFIG).(*gorm.DB)",

		ServiceStr:     "",   // 一般不填 service层名称 一般只有 2 个返回 data,err
		ReturnDataBool: true, // 是否需要返回数据

		LogOrSave: "", // 默认不填
	}
	fmt.Println(GenerateController(c))
}

func TestModels(t *testing.T) {
	t.Run("models", func(t *testing.T) {
		normal := Normal{
			DataSource: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "10.0.10.209", "5432", "ytf", "ytf", "disable", "ytf@2021"),
			TableName:  "ytf_adv_account_info",
			Driver:     "postgres",
		}

		generate, err := AutoGenerateModel(&normal)
		assert.Nil(t, err)

		fmt.Printf(generate)
	})
}

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

	crud := AutoGenerateCRUD(st)
	fmt.Printf(crud)
}
