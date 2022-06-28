package auto

import (
	"reflect"
	"strings"
)

type CRUD interface {
	initStruct()
	formatCRUD() string
	formatSqlxCRUD() string
}

type St struct {
	Stru interface{}            `json:"stru"` // 结构体
	Info map[string]interface{} `json:"info"` // 无关痛痒，缺了也没关系

	DbConfig    string `json:"db_config"`    // 配置数据库
	RedisConfig string `json:"redis_config"` // 配置缓存
	ModelsName  string `json:"models_name"`  // 配置层级
	LogOrSave   string `json:"log_or_save"`  // 配置 打印 默认  log.Printf("%v",errors.WithStack(err))
	Handlers    string `json:"handlers"`     // 所在包名

	valueNameList []string `json:"value_name_list"` // 结构体内部名称
}

//  AutoGenerateCRUD 自动生成结构体名称
func AutoGenerateCRUD(c CRUD) string {

	c.initStruct()

	crud := c.formatCRUD()

	return crud
}

func (s *St) initStruct() {
	s.parseName()

	if s.LogOrSave == "" {
		s.LogOrSave = `log.Printf("%+v",errors.WithStack(err))`
	}
	if _, ok := s.Info["tag"]; !ok {
		s.Info["tag"] = "测试"
	}
	s.parseValue()

}

// 将名称大小写分组
func (s *St) parseName() {
	if s.Info == nil {
		s.Info = map[string]interface{}{}
	}
	// 分组
	name := reflect.TypeOf(s.Stru).Name()
	s.Info["name"] = name
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

	s.Info["name_list"] = nameList
}

// value name
func (s *St) parseValue() {
	typeOf := reflect.TypeOf(s.Stru)

	valeNameList := make([]string, 0, typeOf.NumField())
	for i := 0; i < typeOf.NumField(); i++ {
		valeNameList = append(valeNameList, typeOf.Field(i).Name)
	}
	s.valueNameList = valeNameList
}

func (s St) formatCRUD() string {
	var tmp string

	tmp = s.GenerateBefore(tmp)
	tmp = s.GenerateAdd(tmp)
	tmp = s.GenerateList(tmp)
	tmp = s.GenerateGet(tmp)
	tmp = s.GenerateUpdate(tmp)
	tmp = s.GenerateDelete(tmp)
	return tmp
}

// ${t_name} 横线名称
// ${name} 开头大写名称
// ${l_name} 小名称
// ${models_name} 层级名称

func (s St) GenerateBefore(tmp string) string {

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
	str = strings.ReplaceAll(str, "${t_name}", s.parseLName())
	str = strings.ReplaceAll(str, "${l_name}", s.Handlers)
	str = strings.ReplaceAll(str, "${name}", s.Info["name"].(string))
	return tmp + "\n" + str
}

// 转换为接口名称
func (s St) parseLName() string {

	var tmp strings.Builder
	for _, str := range s.Info["name_list"].([]string) {
		tmp.WriteString(strings.ToLower(str))
		tmp.WriteString("_")
	}
	strl := tmp.String()
	return strl[:len(strl)-1]
}

// ${models_name} 结构体
// ${name} 结构体名称
// ${log_or_save} 自定义输入输出
// ${redis_config} 缓存配置
// ${db_config}  数据库配置

func (s St) GenerateAdd(tmp string) string {
	var str = `
// Auto generate add
func HttpAdd${name}(c *gin.Context) {
	var param ${models_name}
	if err := c.Bind(&param); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	// need rely
	if err := (${db_config}).Model(&param).Create(&param).Error; err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 2, "msg": err.Error()})
		return
	}

	if param.RedisKey() != "" {
		// need redis
		conn := (${redis_config}).(redis.Conn)
		defer conn.Close()
		_ = param.DeleteFromRedis(conn)
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success", "data": param})
}
`
	str = strings.ReplaceAll(str, "${models_name}", s.ModelsName)
	str = strings.ReplaceAll(str, "${name}", s.Info["name"].(string))
	str = strings.ReplaceAll(str, "${log_or_save}", s.LogOrSave)
	str = strings.ReplaceAll(str, "${redis_config}", s.RedisConfig)
	str = strings.ReplaceAll(str, "${db_config}", s.DbConfig)

	return tmp + "\n" + str
}

// ${models_name} 结构体
// ${name} 结构体名称
// ${log_or_save} 自定义输入输出
// ${db_config}  数据库配置
// ${auto_where} 自动生成条件

func (s St) GenerateList(tmp string) string {
	var str = `
// auto generate list get
func HttpList${name}(c *gin.Context) {
	db := (${db_config}).Model(${models_name}{})

	createdAtStartTimeStr := c.DefaultQuery("created_at_start", "")
	createdAtEndTimeStr := c.DefaultQuery("created_at_end", "")
	var createdAtStart, createdAtEnd time.Time
	if createdAtStartTimeStr != "" && createdAtEndTimeStr != "" {
		var e error
		createdAtStart, e = time.ParseInLocation("2006-01-02", createdAtStartTimeStr, time.Local)
		if e != nil {
			c.JSON(400, gin.H{"code": 1, "msg": e.Error()})
			return
		}
		createdAtEnd, e = time.ParseInLocation("2006-01-02", createdAtEndTimeStr, time.Local)
		if e != nil {
			c.JSON(400, gin.H{"code": 2, "msg": e.Error()})
			return
		}
		db = db.Where("created_at between ? and ?", createdAtStart, createdAtEnd.AddDate(0, 0, 1))
	}
	updatedAtStartTimeStr := c.DefaultQuery("updated_at_start", "")
	updatedAtEndTimeStr := c.DefaultQuery("updated_at_end", "")
	var updatedAtStart, updatedAtEnd time.Time
	if updatedAtStartTimeStr != "" && updatedAtEndTimeStr != "" {
		var e error
		updatedAtStart, e = time.ParseInLocation("2006-01-02", updatedAtStartTimeStr, time.Local)
		if e != nil {
			c.JSON(400, gin.H{"code": 3, "msg": e.Error()})
			return
		}
		updatedAtEnd, e = time.ParseInLocation("2006-01-02", updatedAtEndTimeStr, time.Local)
		if e != nil {
			c.JSON(400, gin.H{"code": 4, "msg": e.Error()})
			return
		}
		db = db.Where("updated_at between ? and ?", updatedAtStart, updatedAtEnd.AddDate(0, 0, 1))
	}

${auto_where}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		${log_or_save}
		c.JSON(500, gin.H{"code": 5, "msg": err.Error()})
		return
	}
	list := make([]${models_name}, 0, 20)
	if count == 0 {
		c.JSON(200, gin.H{"code": 0, "msg": "success", "count": count, "data": list})
		return
	}

	// limit
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "20")
	limit, offset := utilx.ToLimitOffset(size, page, count)
	db = db.Limit(limit).Offset(offset)

	// order by
	if orderBy := c.DefaultQuery("order_by", ""); orderBy != "" {
		db.Order(utilx.ParseOrderBy(orderBy))
	}

	// Find
	if err := db.Find(&list).Error; err != nil {
		${log_or_save}
		c.JSON(500, gin.H{"code": 6, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success", "count": count, "data": list})
}
`
	str = strings.ReplaceAll(str, "${models_name}", s.ModelsName)
	str = strings.ReplaceAll(str, "${name}", s.Info["name"].(string))
	str = strings.ReplaceAll(str, "${log_or_save}", s.LogOrSave)
	str = strings.ReplaceAll(str, "${db_config}", s.DbConfig)
	str = strings.ReplaceAll(str, "${auto_where}", s.autoWhere())

	return tmp + "\n" + str
}

// ${models_name}
// ${value_name}
// ${tmp_name}

// 支持 json 传参
func (s St) autoWhereTwo() string {
	var tmp = `
	var param ${models_name}
	if err := c.Bind(&param); err != nil {
		log.Printf("%+v", errors.WithStack(err))
		c.JSON(200, gin.H{"code": 5, "msg": err.Error()})
		return
	}
`

	for _, name := range s.valueNameList {
		var str = `
	if len(param.${value_name}) >0 {
		db = db.Where("${tmp_name}=?", ${value_name})
	}
`
		str = strings.ReplaceAll(str, "${value_name}", name)
		str = strings.ReplaceAll(str, "${tmp_name}", parseL(name))
		tmp += str
	}

	return tmp
}

// ${value_name}
// ${tmp_name}

// 自动生成条件
func (s St) autoWhere() string {
	var tmp strings.Builder
	for _, name := range s.valueNameList {
		var str = `
	if ${value_name} := c.DefaultQuery("${tmp_name}", ""); ${value_name} != "" {
		db = db.Where("${tmp_name}=?", ${value_name})
	}
`
		str = strings.ReplaceAll(str, "${value_name}", name)
		str = strings.ReplaceAll(str, "${tmp_name}", parseL(name))
		tmp.WriteString(str)
	}
	return tmp.String()
}

// 转下划线
func parseL(name string) string {
	var tmp strings.Builder
	var index, count int
	for i, str := range strings.Split(name, "") {
		if i == 0 {
			continue
		}
		if str >= "A" && str <= "Z" {
			tmp.WriteString(strings.ToLower(name[index:i]))
			tmp.WriteString("_")
			index = i
			count++
		}
	}

	tmp.WriteString(name[index:])

	tmps := tmp.String()

	return strings.ToLower(tmps)
}

// ${models_name} 结构体
// ${name} 结构体名称
// ${log_or_save} 自定义输入输出
// ${db_config}  数据库配置

func (s St) GenerateGet(tmp string) string {
	var str = `
// Auto generate get
func HttpGet${name}(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"code": 1, "msg": fmt.Sprintf("param 'id' requires int but got %d", id)})
		return
	}

	var count int64
	if err = (${db_config}).Model(${models_name}{}).Where("id=?", id).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"code": 2, "msg": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(200, gin.H{"code": 3, "msg": fmt.Sprintf("id '%s' data not found", c.Param("id"))})
		return
	}

	var resp ${models_name}
	if err = (${db_config}).Model(${models_name}{}).Where("id=?", id).First(&resp).Error; err != nil {
		${log_or_save}
		c.JSON(500, gin.H{"code": 4, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success", "data": resp})
}
`
	str = strings.ReplaceAll(str, "${models_name}", s.ModelsName)
	str = strings.ReplaceAll(str, "${name}", s.Info["name"].(string))
	str = strings.ReplaceAll(str, "${log_or_save}", s.LogOrSave)
	str = strings.ReplaceAll(str, "${db_config}", s.DbConfig)

	return tmp + "\n" + str
}

// ${models_name} 结构体
// ${name} 结构体名称
// ${log_or_save} 自定义输入输出
// ${db_config}  数据库配置
// ${redis_config}

func (s St) GenerateUpdate(tmp string) string {

	var str = `
// Auto generate update
func HttpUpdate${name}(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"code": 1, "msg": fmt.Sprintf("param 'id' requires int but got %d", id)})
		return
	}
	var count int64
	if err = (${db_config}).Model(${models_name}{}).Where("id=?", id).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"code": 2, "msg": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(200, gin.H{"code": 3, "msg": fmt.Sprintf("id '%s' data not found", c.Param("id"))})
		return
	}
	var param ${models_name}
	if err = c.Bind(&param); err != nil {
		c.JSON(400, gin.H{"code": 4, "msg": err.Error()})
		return
	}

	var resp ${models_name}
	if err = (${db_config}).Model(${models_name}{}).Where("id=?", id).First(&resp).Error; err != nil {
		${log_or_save}
		c.JSON(500, gin.H{"code": 5, "msg": err.Error()})
		return
	}

	tx := (${db_config}).Begin()
	if err = tx.Model(${models_name}{}).Where("id=?", id).Updates(param).Error; err != nil {
		tx.Rollback()
		${log_or_save}
		c.JSON(500, gin.H{"code": 6, "msg": err.Error()})
		return
	}

	tx.Commit()
	if resp.RedisKey() != "" {
		conn := ${redis_config}
		defer conn.Close()
		_ = resp.DeleteFromRedis(conn)
	}
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}
`
	str = strings.ReplaceAll(str, "${models_name}", s.ModelsName)
	str = strings.ReplaceAll(str, "${name}", s.Info["name"].(string))
	str = strings.ReplaceAll(str, "${log_or_save}", s.LogOrSave)
	str = strings.ReplaceAll(str, "${db_config}", s.DbConfig)
	str = strings.ReplaceAll(str, "${redis_config}", s.RedisConfig)
	return tmp + "\n" + str
}

// ${models_name} 结构体
// ${name} 结构体名称
// ${log_or_save} 自定义输入输出
// ${db_config}  数据库配置

func (s St) GenerateDelete(tmp string) string {
	var str = `
// Auto generate delete
func HttpDelete${name}(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"code": 1, "msg": fmt.Sprintf("param 'id' requires int but got %d", id)})
		return
	}

	var count int64
	if err = (${db_config}).Model(${models_name}{}).Where("id=?", id).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"code": 2, "msg": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(200, gin.H{"code": 3, "msg": fmt.Sprintf("id '%s' data not found", id)})
		return
	}

	var resp ${models_name}
	if err = (${db_config}).Model(${models_name}{}).Where("id=?", id).First(&resp).Error; err != nil {
		${log_or_save}
		c.JSON(500, gin.H{"code": 4, "msg": err.Error()})
		return
	}

	if err = (${db_config}).Model(${models_name}{}).Where("id=?", id).Delete(&${models_name}{}).Error; err != nil {
		${log_or_save}
		c.JSON(500, gin.H{"code": 5, "msg": err.Error()})
		return
	}

	if resp.RedisKey() != "" {
		conn := ${redis_config}
		defer conn.Close()
	}
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}
`
	str = strings.ReplaceAll(str, "${models_name}", s.ModelsName)
	str = strings.ReplaceAll(str, "${name}", s.Info["name"].(string))
	str = strings.ReplaceAll(str, "${log_or_save}", s.LogOrSave)
	str = strings.ReplaceAll(str, "${db_config}", s.DbConfig)
	str = strings.ReplaceAll(str, "${redis_config}", s.RedisConfig)

	return tmp + "\n" + str
}
