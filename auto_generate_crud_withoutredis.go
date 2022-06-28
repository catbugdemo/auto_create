package auto

import (
	"fmt"
	"reflect"
	"strings"
)

type WithoutRedis struct {
	Stru interface{}            `json:"stru"` // 结构体
	Info map[string]interface{} `json:"info"` // 重要层级
	//  "controller":"handler" , "service":"service" , "models":"models"

	ToMany ToMany

	DbConfig  string `json:"db_config"`   // 配置数据库
	LogOrSave string `json:"log_or_save"` // 配置打印

	valueNameList []string `json:"value_name_list"` // 结构体内部名称
	commonList    []string `json:"common_list"`
}

type ToMany struct {
	// 多 对 多 表关联
	ConnectTable   interface{} // 多对多连接表
	BeConnectTable interface{} // 被关联表
}

func (o *WithoutRedis) initStruct() {
	if o.Info == nil {
		o.Info = map[string]interface{}{}
	}

	o.parseName()

	if o.LogOrSave == "" {
		o.LogOrSave = `log.Printf("%+v", errors.WithStack(err))`
	}

	o.parseValue()

	if (o.ToMany != ToMany{}) {
		o.parseConnect()
	}
}

// 将名称大小写分组
func (o *WithoutRedis) parseName() {
	if o.Info == nil {
		o.Info = map[string]interface{}{}
	}
	// 分组
	name := reflect.TypeOf(o.Stru).Name()
	o.Info["name"] = name
	o.Info["req_name"] = name + "Req"
	o.Info["resp_name"] = name + "Resp"
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
	o.Info["model_name"] = o.Info["models"].(string) + "." + name
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

// 比较相同，同时进行关联
func (o *WithoutRedis) parseConnect() {
	// 查找关联数据
	BeConnectTable := reflect.TypeOf(o.ToMany.BeConnectTable)

	o.Info["be_connect_table_name"] = BeConnectTable.Name()

	ConnectTable := reflect.TypeOf(o.ToMany.ConnectTable)

	o.Info["connect_table_name"] = ConnectTable.Name()

	list := make([]string, 0, ConnectTable.NumField())
	// 将 主表和关联表找出相同参数
	for i := 0; i < ConnectTable.NumField(); i++ {
		for _, str := range o.valueNameList {
			if str == "CreatedAt" || str == "UpdatedAt" || str == "Id" {
				continue
			}

			// 如果固定相同
			if str == ConnectTable.Field(i).Name {
				list = append(list, str)
			}
		}
	}
	o.commonList = list

}

func (o *WithoutRedis) formatCRUD() string {
	var tmp string
	tmp = o.GenerateBefore(tmp)
	tmp = o.GenerateController(tmp)
	tmp = o.GenerateService(tmp)
	return tmp
}

func (o WithoutRedis) GenerateBefore(tmp string) string {

	var str = `
import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/catbugdemo/errors" 	// or "github.com/pkg/errors"

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
	str = strings.ReplaceAll(str, "${l_name}", o.Info["controller"].(string))
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
// ${model_name} 底层体名称
// ${service} 请求地址
// ${log_or_save}

func (o WithoutRedis) GenerateController(tmp string) string {
	var str = `
// auto generate add
func HttpAdd${name}(c *gin.Context) {
	var req ${service}.${req_name}
	if err := c.Bind(&req); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 1, "msg": "HttpAdd${name} request binding failed", "debug": err.Error()})
		return
	}
	data, err := ${service}.HttpAdd${name}(req, ${db_config}); 
	if err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 2, "msg": "HttpAdd${name} service operate failed", "debug": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success", "data": data})
}

// auto generate get 
func HttpGet${name}(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		${log_or_save}
		c.JSON(400, gin.H{"code": 1, "msg": "HttpGet${name} request binding failed", "debug": fmt.Sprintf("param 'id' requires int but got %+v", id)})
		return
	}

	data, err :=${service}.HttpGet${name}(id, ${db_config})
	if err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 2, "msg": "HttpGet${name} service operate failed", "debug": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success", "data": data})
}


// auto generate update
func HttpUpdate${name}(c *gin.Context) {
	var req ${service}.${req_name}
	if err := c.Bind(&req); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 1, "msg": "HttpUpdate${name} request binding failed", "debug": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		${log_or_save}
		c.JSON(400, gin.H{"code": 2, "msg": fmt.Sprintf("param 'id' requires int but got %v", id)})
		return
	}
	
	if err = ${service}.HttpUpdate${name}(req, ${db_config}); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 3, "msg": "HttpUpdate${name} service operate failed", "debug": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

// auto generate delete 
func HttpDelete${name}(c *gin.Context) {
	var req ${service}.${req_name}
	if err := c.Bind(&req); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 1, "msg": "HttpDelete${name} request binding failed", "debug": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 2, "msg": "HttpDelete${name} service operate failed", "debug": err.Error()})
		return
	}
	
	if err = ${service}.HttpDelete${name}(id, ${db_config}); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 3, "msg": "service operate failed", "debug": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

// auto generate get list
func HttpList${name}(c *gin.Context) {	
	eng := (${db_config}).Model(${model_name}{})
	List${name}Condition(c,eng)

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "20")
	orderBy := c.DefaultQuery("order_by", "")

	count, list, err := ${service}.HttpList${name}(page, size, orderBy, eng)
	if err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 1, "msg": "HttpList${name} service operate failed", "debug": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"code": 0, "msg": "success", "count": count, "data": list})
}

// List condition
func List${name}Condition(c *gin.Context,db *gorm.DB) {
	${auto_where}
}
`
	str = strings.ReplaceAll(str, "${name}", o.Info["name"].(string))
	str = strings.ReplaceAll(str, "${req_name}", o.Info["req_name"].(string))
	str = strings.ReplaceAll(str, "${model_name}", o.Info["model_name"].(string))
	str = strings.ReplaceAll(str, "${service}", o.Info["service"].(string))
	str = strings.ReplaceAll(str, "${log_or_save}", o.LogOrSave)
	str = strings.ReplaceAll(str, "${auto_where}", o.autoWhere())
	str = strings.ReplaceAll(str, "${db_config}", o.DbConfig)

	return tmp + str
}

// 自动生成条件
func (o WithoutRedis) autoWhere() string {
	var tmp strings.Builder
	for _, name := range o.valueNameList {
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

// ${name} 名称
// ${req_name} 请求参数
// ${model_name} 结构体位置
// ${connect_ids} 关联 id 组 支持多关联组
// ${db_config} 配置
// ${connect_add_operate} 关联操作

func (o WithoutRedis) GenerateService(tmp string) string {
	var str = `
// ------------------------------------------------ auto service crud

type ${req_name} struct {
	${model_name}
	${connect_ids} 
}

type ${resp_name} struct {
	${model_name}
	${be_connect_table_list} 
}

// auto generate service add
func HttpAdd${name}(req ${req_name},db *gorm.DB) (${req_name}, error) {
	param := req.${name}

	// operate basic data
	tx := db.Begin()
	defer tx.Commit()

	if err := tx.Model(&${model_name}{}).Create(&param).Error; err != nil {
		return errors.WithStack(err)
	}

${connect_add_operate}
	return nil
}

// auto generate service get 
func HttpGet${name}(id int,db *gorm.DB)(${resp_name}, error) {
	param, err := ${name}FirstIfExist(id,db)
	if err != nil {
		return ${resp_name}{}, errors.WithStack(err)
	}

${connect_get_operate}	
	return ${resp_name}{
		${name}: param,
		${connect_get_data}
	}, nil
}

${connect_get_func}

// auto generate service update
func HttpUpdate${name}(req ${req_name},db *gorm.DB) error {
	param, err := ${name}FirstIfExist(req.Id, db)
	if err != nil {
		return errors.WithStack(err)
	}

	// --- update
	tx := db.Begin()
	defer tx.Commit()

	if err = tx.Model(${model_name}{}).Where("id=?", param.Id).Updates(param).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

${connect_update}
	return nil
}
${connect_update_func}

// auto generate service delete
func HttpDelete${name}(id int, db *gorm.DB) error {
	param, err := ${name}FirstIfExist(id,db)
	if err != nil {
		return errors.WithStack(err)
	}

	// --- delete
	tx := db.Begin()
	defer tx.Commit()

	if err = tx.Model(${model_name}{}).Where("id=?", id).Delete(&${model_name}{}).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}
${connect_delete}
	return nil// Find
}


// auto generate service list
func HttpList${name}(page string, size string, orderBy string,eng *gorm.DB, db *gorm.DB) (int, []${resp_name}, error) {
	var count int64
	if err := eng.Count(&count).Error; err != nil {
		return 0, nil,  errors.WithStack(err)
	}
	list := make([]${model_name}, 0, count)
	if count == 0 {
		return 0, nil, nil
	}
	limit, offset := utils.ToLimitOffset(size, page, int(count))
	db = db.Limit(limit).Offset(offset)
	if orderBy != "" {
		eng = eng.Order(utils.GenerateOrderBy(orderBy))
	}
	// Find
	if err := db.Find(&list).Error; err != nil {
		return	0, nil, errors.WithStack(err)
	}	
	// ------- 
${connect_list}
	return int(count), respList, nil
}


// FirstIfExist
func ${name}FirstIfExist(id int, db *gorm.DB) (${model_name}, error) {
	var count int64
	if err := db.Model(${model_name}{}).Where("id=?", id).Count(&count).Error; err != nil {
		return ${model_name}{}, errors.WithStack(err)
	}

	if count == 0 {
		return ${model_name}{}, errors.New(fmt.Sprintf("id '%d' data not found", id))
	}

	// first
	var param ${model_name}
	if err := db.Model(${model_name}{}).Where("id=?", id).First(&param).Error; err != nil {
		return ${model_name}{}, errors.WithStack(err)
	}
	
	return param, nil
}
`

	// relation
	str = strings.ReplaceAll(str, "${connect_ids}", o.checkoutConnectIds())
	str = strings.ReplaceAll(str, "${connect_add_operate}", o.printAddConnect())
	str = strings.ReplaceAll(str, "${connect_get_operate}", o.printGetOperate())
	str = strings.ReplaceAll(str, "${connect_get_data}", o.printGetData())
	str = strings.ReplaceAll(str, "${connect_get_func}", o.printGetConnect())
	// update
	str = strings.ReplaceAll(str, "${connect_update}", o.printGetUpdate())
	str = strings.ReplaceAll(str, "${connect_update_func}", o.printUpdateConnectfunc())
	str = strings.ReplaceAll(str, "${connect_delete}", o.printDeleteConnect())
	str = strings.ReplaceAll(str, "${connect_list}", o.printGetList())
	str = strings.ReplaceAll(str, "${be_connect_table_list}", o.printTableList())

	// replace
	str = strings.ReplaceAll(str, "${be_connect_table_name}", o.Info["be_connect_table_name"].(string))
	str = strings.ReplaceAll(str, "${connect_table_name}", o.Info["connect_table_name"].(string))
	str = strings.ReplaceAll(str, "${resp_name}", o.Info["resp_name"].(string))

	str = strings.ReplaceAll(str, "${name}", o.Info["name"].(string))
	str = strings.ReplaceAll(str, "${req_name}", o.Info["req_name"].(string))
	str = strings.ReplaceAll(str, "${model_name}", o.Info["model_name"].(string))
	str = strings.ReplaceAll(str, "${service}", o.Info["service"].(string))
	str = strings.ReplaceAll(str, "${log_or_save}", o.LogOrSave)
	str = strings.ReplaceAll(str, "${db_config}", o.DbConfig)

	return tmp + str
}

func (o WithoutRedis) printTableList() string {
	if (o.ToMany == ToMany{}) {
		return ""
	}
	var str = `${be_connect_table_name}s []${models}.${be_connect_table_name}`
	return str
}

func (o WithoutRedis) printGetList() string {
	var str = `
	respList := make([]${resp_name} ,0 , count)
	for _, param := range list {
		l ,err := GetList${be_connect_table_name}ById(param, db)
		if err != nil {
			log.Printf("%+v",errors.WithStack(err))
		}
		respList = append(respList, ${resp_name}{
			${name}: param,
			${resp_name}s
		})
	}
`

	return str
}

func (o WithoutRedis) printDeleteConnect() string {
	if (ToMany{} == o.ToMany) {
		return ""
	}
	var str = `
	// delete connect table
	if err := tx.Model(${models}.${connect_table_name}{}).Where(${where}).Delete(&${models}.${connect_table_name}{}).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}
`
	return str
}

func (o WithoutRedis) printGetUpdate() string {
	if (o.ToMany == ToMany{}) {
		return ""
	}

	var str = `
	if err := Update${connect_table_name}(req, tx);err!=nil {
		tx.Rollback()
		return errors.WithStack(err)
	}
`
	return str
}

func (o WithoutRedis) checkoutConnectIds() string {
	if (o.ToMany == ToMany{}) {
		return ""
	}

	return fmt.Sprintf("%vIds []int //TODO json", o.Info["be_connect_table_name"])
}

func (o WithoutRedis) printGetOperate() string {
	if (ToMany{} == o.ToMany) {
		return ""
	}

	var str = `
	list ,err := GetList${be_connect_table_name}ById(param, db)
	if err != nil {
		return ${resp_name}{}, errors.WithStack(err)
	}
`
	return str
}

func (o WithoutRedis) printGetData() string {
	if (ToMany{} == o.ToMany) {
		return ""
	}
	return "${be_connect_table_name}s: list,"
}

// ${be_connect_table_name}
// ${models}
// ${connect_table_name}
// ${fix} 相同参数

func (o WithoutRedis) printAddConnect() string {
	if (ToMany{} == o.ToMany) {
		return ""
	}

	var str = `
	// operate connect table
	if len(req.${be_connect_table_name}Ids) == 0 {
		return nil
	}
	
	connect := make([]${models}.${connect_table_name}, 0, len(req.${be_connect_table_name}Ids))
	for _, id := range req.${be_connect_table_name}Ids {
		connect = append(connect, models.${connect_table_name}{
			${fix}
			${name}Id: param.Id,     
			${be_connect_table_name}Id:   id,          
		})
	}

	if err := tx.Model(&${models}.${connect_table_name}{}).Create(&connect).Error; err != nil {
		return errors.WithStack(err)
	}
`

	str = strings.ReplaceAll(str, "${be_connect_table_name}", o.Info["be_connect_table_name"].(string))
	str = strings.ReplaceAll(str, "${connect_table_name}", o.Info["connect_table_name"].(string))
	str = strings.ReplaceAll(str, "${models}", o.Info["models"].(string))
	str = strings.ReplaceAll(str, "${fix}", o.dealFix())
	return str
}

func (o WithoutRedis) dealFix() string {
	var tmp string
	for _, comm := range o.commonList {
		tmp += fmt.Sprintf("%v : param.%v,\n", comm, comm)
	}
	return tmp
}

// 作为一个结构体自动生成
func (o WithoutRedis) printGetConnect() string {
	if (ToMany{} == o.ToMany) {
		return ""
	}

	var str = `
func GetList${be_connect_table_name}ById(param ${model_name}, db *gorm.DB) ([]${models}.${be_connect_table_name}, error) {
	// connect table
	type Connect struct {
		${be_connect_table_name}Id []int // TODO json 
	}
	var connect Connect 

	// todo
	if err := db.Model(&${models}.${connect_table_name}{}).Select("${l_be_name}").Where("${where}",${param}).Find(&connect.${be_connect_table_name}).Error;err!=nil{
		return nil, errors.WithStack(err)
	}

	if len(connect.${be_connect_table_name}Id) == 0 {
		return nil, nil
	}

	// select be connect table
	list := make([]${models}.${be_connect_table_name},0,len(connect.${be_connect_table_name}Id))
	if err := db.Model(&${models}.${be_connect_table_name}{}).Where("id in (?)", connect.${be_connect_table_name}Id).Find(&list).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}
`
	str = strings.ReplaceAll(str, "${l_be_name}_id", LBeName(o.Info["be_connect_table_name"].(string)))
	str = strings.ReplaceAll(str, "${where}", o.printWhere())
	str = strings.ReplaceAll(str, "${param}", o.printParam())
	return str
}

func (o WithoutRedis) printParam() string {
	var str = ``
	for _, s := range o.commonList {
		str += "param." + s + "and "
	}
	str += "param.Id"
	return str
}

func (o WithoutRedis) printWhere() string {
	var str = ``
	// 处理县共同数据
	for _, s := range o.commonList {
		str += LBeName(s) + `=? and`
	}

	// 处理 Id 数据
	str += fmt.Sprintf(" %v_id=?", LBeName(o.Info["name"].(string)))
	return str
}

func LBeName(name string) string {
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

	var tmp strings.Builder

	for _, str := range nameList {
		tmp.WriteString(strings.ToLower(str))
		tmp.WriteString("_")
	}
	strl := tmp.String()
	return strl[:len(strl)-1]
}

// 修改操作
func (o WithoutRedis) printUpdateConnectfunc() string {
	if (ToMany{} == o.ToMany) {
		return ""
	}

	var str = `
func Update${connect_table_name}(req ${req_name}, tx *gorm.DB) error {
	if len(req.${be_connect})


	if err := tx.Model(${models}.${connect_table_name}{}).Where("${where}",${param}).Delete(&models.YtfAdvIndustryPeople{}).Error; err != nil {
		return errors.WithStack(err)
	}

	// 处理关联数据
	connect := make([]${models}.${connect_table_name}, 0, len(req.${connect_table_name}Ids))
	for _, id := range req.YtfAdvPeopleTargetIds {
		connect = append(connect, models.YtfAdvIndustryPeople{
			${fix}
			${name}Id: param.Id,     
			${be_connect_table_name}Id:   id,  
		})
	}

	// 添加关联表
	if err := tx.Model(&models.YtfAdvIndustryPeople{}).Create(&connect).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil 
}
`
	str = strings.ReplaceAll(str, "${where}", o.printWhere())
	str = strings.ReplaceAll(str, "${param}", o.printParam())
	str = strings.ReplaceAll(str, "${fix}", o.dealFix())

	return str
}
