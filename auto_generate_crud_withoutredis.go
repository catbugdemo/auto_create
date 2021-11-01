package auto

import (
	"reflect"
	"strings"
)

type Relationship string

const (
	ONE_TO_ONE   Relationship = "ONE_TO_ONE"
	ONE_TO_MANY  Relationship = "ONE_TO_MANY"
	MANY_TO_MANY Relationship = "MANY_TO_MANY"
)

type WithoutRedis struct {
	Stru interface{}            `json:"stru"` // 结构体
	Info map[string]interface{} `json:"info"` // 重要层级
	//  "controller":"handler" , "service":"service" , "model":"models"

	// 多 对 多 表关联
	Relationship Relationship // 关联形式

	DbConfig  string `json:"db_config"`   // 配置数据库
	LogOrSave string `json:"log_or_save"` // 配置打印
	Handlers  string `json:"handlers"`    // 所在包名

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
	var tmp string
	tmp = o.GenerateBefore(tmp)
	//o.GenerateController()
	return ""
}

func (o WithoutRedis) GenerateBefore(tmp string) string {

	var str = `
import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/catbugdemo/errors" 	// or "github.com/pkg/errors"

	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

/*
	Auto generate crud
	need pkg
		- "github.com/gin-gonic/gin"
		- "gorm.io/gorm"
		- "github.com/catbugdemo/errors" or "github.com/pkg/errors"
		- "github.com/gomodule/redigo/redis"
		- "github.com/catbugdemo/utilx"

	when generate can use gin :
		r := gin.Default
		r.POST("/${t_name}/",${l_name}.HttpAdd${name})
		r.PATCH("/${t_name}/:id/",${l_name}.HttpUpdate${name})
		r.DELETE("/${t_name}/:id/",${l_name}.HttpDelete${name})
		r.GET("/${t_name}/",${l_name}.HttpList${name})
		r.GET("/${t_name}/:id/",${l_name}.HttpGet${name})
*/

`
	str = strings.ReplaceAll(str, "${t_name}", o.parseLName())
	str = strings.ReplaceAll(str, "${l_name}", o.Handlers)
	str = strings.ReplaceAll(str, "${name}", o.Info["name"].(string))
	return tmp + "\n" + str
}

func (o WithoutRedis) parseLName() string {

	var tmp strings.Builder
	for _, str := range o.Info["name_list"].([]string) {
		tmp.WriteString(strings.ToLower(str))
		tmp.WriteString("_")
	}
	strl := tmp.String()
	return strl[:len(strl)-1]
}

// ${name} 名称
// ${req_name} req 名称

/*func (o WithoutRedis) GenerateController(tmp string) string {
	var str = `
// auto generate add
func HttpAdd${name}(c *gin.Context) {
	var req ${req_name}
	if err := c.Bind(&req); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 1, "msg": "HttpAdd${name} request binding failed", "debug": err.Error()})
		return
	}
	db := ${db_config}
	if err := service.HttpAdd${name}(&req, db); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 2, "msg": "HttpAdd${name} service operate failed", "debug": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success", "data": req})
}

// auto generate get list

`

	return ""
}
*/

/*
func (o WithoutRedis) GenerateService() string {
	var str = `
// ------------------------------------------------ auto service crud

type ${req_name} struct {

}


// auto generate service add

`
	return ""
}



return ""
}
*/
