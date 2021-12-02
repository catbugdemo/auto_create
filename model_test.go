package auto

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

func TestAGenerateCRUD(t *testing.T) {
	type YtfAdvPeopleTarget struct {
		Id        int       `gorm:"column:id;default:" json:"id" form:"id"`
		CreatedAt time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`

		UserId        int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		CreatedPeople string `gorm:"column:created_people;default:" json:"created_people" form:"created_people"`
		UpdatedPeople string `gorm:"column:updated_people;default:" json:"updated_people" form:"updated_people"`

		PeopleName                string         `gorm:"column:people_name;default:" json:"people_name" form:"people_name"`
		Age                       postgres.Jsonb `gorm:"column:age;default:" json:"age" form:"age"`
		Gender                    string         `gorm:"column:gender;default:" json:"gender" form:"gender"`
		Education                 string         `gorm:"column:education;default:" json:"education" form:"education"`
		MaritalStatus             string         `gorm:"column:marital_status;default:" json:"marital_status" form:"marital_status"`
		WorkingStatus             string         `gorm:"column:working_status;default:" json:"working_status" form:"working_status"`
		FinancialSituation        string         `gorm:"column:financial_situation;default:" json:"financial_situation" form:"financial_situation"`
		ConsumptionType           string         `gorm:"column:consumption_type;default:" json:"consumption_type" form:"consumption_type"`
		GameConsumptionLevel      string         `gorm:"column:game_consumption_level;default:" json:"game_consumption_level" form:"game_consumption_level"`
		ConsumptionStatus         string         `gorm:"column:consumption_status;default:" json:"consumption_status" form:"consumption_status"`
		ResidentialCommunityPrice postgres.Jsonb `gorm:"column:residential_community_price;default:" json:"residential_community_price" form:"residential_community_price"`
		BehaviorOrInterest        postgres.Jsonb `gorm:"column:behavior_or_interest;default:" json:"behavior_or_interest" form:"behavior_or_interest"`
		NewDevice                 string         `gorm:"column:new_device;default:" json:"new_device" form:"new_device"`
		ExcludedConvertedAudience postgres.Jsonb `gorm:"column:excluded_converted_audience;default:" json:"excluded_converted_audience" form:"excluded_converted_audience"`
		DeprecatedCustomAudience  postgres.Jsonb `gorm:"column:deprecated_custom_audience;default:" json:"deprecated_custom_audience" form:"deprecated_custom_audience"`
		DeviceBrandModel          postgres.Jsonb `gorm:"column:device_brand_model;default:" json:"device_brand_model" form:"device_brand_model"`
		NetworkScene              string         `gorm:"column:network_scene;default:" json:"network_scene" form:"network_scene"`
		UserOs                    string         `gorm:"column:user_os;default:" json:"user_os" form:"user_os"`
		NetworkType               string         `gorm:"column:network_type;default:" json:"network_type" form:"network_type"`
		NetworkOperator           string         `gorm:"column:network_operator;default:" json:"network_operator" form:"network_operator"`
		DevicePrice               string         `gorm:"column:device_price;default:" json:"device_price" form:"device_price"`
		MobileUnionCategory       string         `gorm:"column:mobile_union_category;default:" json:"mobile_union_category" form:"mobile_union_category"`
		Temperature               string         `gorm:"column:temperature;default:" json:"temperature" form:"temperature"`
		UvIndex                   string         `gorm:"column:uv_index;default:" json:"uv_index" form:"uv_index"`
		DressingIndex             string         `gorm:"column:dressing_index;default:" json:"dressing_index" form:"dressing_index"`
		MakeupIndex               string         `gorm:"column:makeup_index;default:" json:"makeup_index" form:"makeup_index"`
		Climate                   string         `gorm:"column:climate;default:" json:"climate" form:"climate"`
		WechatAdBehavior          postgres.Jsonb `gorm:"column:wechat_ad_behavior;default" json:"wechat_ad_behavior" form:"wechat_ad_behavior"`
	}

	type YtfAdvIndustryPeople struct {
		Id                   int       `gorm:"column:id;default:" json:"id" form:"id"`
		CreatedAt            time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		UpdatedAt            time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`
		YtfAdvUserId         int       `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		IndustryId           int       `gorm:"column:industry_id;default:" json:"industry_id" form:"industry_id"`
		YtfAdvPeopleTargetId int       `gorm:"column:people_id;default:" json:"people_id" form:"people_id"`
	}
	type YtfAdvIndustryTemplate struct {
		Id                int       `gorm:"column:id;default:" json:"id" form:"id"`
		CreatedAt         time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		UpdatedAt         time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`
		UserId            int       `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		TemplateName      string    `gorm:"column:template_name;default:" json:"template_name" form:"template_name"`
		CreatedPeople     string    `gorm:"column:created_people;default:" json:"created_people" form:"created_people"`
		UpdatedPeople     string    `gorm:"column:updated_people;default:" json:"updated_people" form:"updated_people"`
		Status            int       `gorm:"column:status;default:" json:"status" form:"status"`
		PrimaryIndustry   int       `gorm:"column:primary_industry;default:" json:"primary_industry" form:"primary_industry"`
		PrimaryName       string    `gorm:"column:primary_name;default:" json:"primary_name" form:"primary_name"`
		SecondaryIndustry int       `gorm:"column:secondary_industry;default:" json:"secondary_industry" form:"secondary_industry"`
		SecondaryName     string    `gorm:"column:secondary_name;default:" json:"secondary_name" form:"secondary_name"`
	}

	wr := WithoutRedis{
		Stru: YtfAdvPeopleTarget{},
		ToMany: ToMany{
			ConnectTable:   YtfAdvIndustryPeople{},
			BeConnectTable: YtfAdvIndustryTemplate{},
		},
		Info: map[string]interface{}{
			"controller": "handlers",
			"service":    "service",
			"model":      "models",
		},

		DbConfig: "c.MustGet(DB_CONFIG).(*gorm.DB)",
	}

	fmt.Println(AutoGenerateCRUD(&wr))
}

func TestSwagger(t *testing.T) {
	swag := SwaggerInit{
		Name:   "素材库",
		Tags:   "ytf_adv_material_library",
		Router: "ytf-adv-material-library",
		Req:    "models.YtfAdvMaterialLibrary",
	}

	swagger := GenerateCURDSwagger(swag)
	fmt.Println(swagger)
}

func TestGenerateNormalSwagger(t *testing.T) {
	normal := SwaggerNormal{
		Describe:      "获取线索配置",                // 描述
		TableName:     "ytf_wxapp_clue_config", //根据表名分组
		IfHeaderToken: true,                    // 是否有安全校验
		Req:           "HttpGetClueConfig",     // 请求参数
		RouterUrl:     "/ytf-wxapp-clue-config",
	}

	swagger := GenerateNormalSwagger(normal)
	fmt.Println(swagger)
}
