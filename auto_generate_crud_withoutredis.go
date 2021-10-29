package auto

import (
	"reflect"
	"strings"
)

type WithoutRedis struct {
	Stru interface{}            `json:"stru"` // 结构体
	Info map[string]interface{} `json:"info"` // 重要层级
	//  "controller":"handler" , "service":"service" , "model":"models"

	DbConfig  string `json:"db_config"`   // 配置数据库
	LogOrSave string `json:"log_or_save"` // 配置打印

	valueNameList []string `json:"value_name_list"` // 结构体内部名称
}

func (o *WithoutRedis) initStruct() {
	o.parseName()

}

// 将名称大小写分组
func (o *WithoutRedis) parseName() {
	if o.Info == nil {
		o.Info = map[string]interface{}{}
	}
	// 分组
	name := reflect.TypeOf(o.Stru).Name()
	o.Info["name"] = name
	var index, count int
	nameList := make([]string, 0, 27)
	for i, str := range strings.Split(name, "") {
		if i == 0 {
			continue
		}

		if str >= "A" && str <= "Z" {
			nameList = append(nameList, name[index:i])
			index = i
			count++
		}
	}

	nameList = append(nameList, name[index:])

	o.Info["name_list"] = nameList
}

// 解析获取 结构体中的数据的名称
func (o *WithoutRedis) parseValue() {
	typeOf := reflect.TypeOf(o.Stru)

	valeNameList := make([]string, 0, typeOf.NumField())
	for i := 0; i < typeOf.NumField(); i++ {
		valeNameList = append(valeNameList, typeOf.Field(i).Name)
	}
	o.valueNameList = valeNameList
}

func (o *WithoutRedis) formatCRUD() string {
	return ""
}

func (o WithoutRedis) GenerateController() string {
	return ""
}

func (o WithoutRedis) GenerateService() string {
	return ""
}

func (o WithoutRedis) GenerateModel() string {
	return ""
}
