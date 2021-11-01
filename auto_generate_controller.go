package auto

import (
	"fmt"
	"strings"
)

type Control struct {
	ControlName string // control名称
	Describe    string // 描述

	Req Req // 请求

	ServiceStr     string // service 不填默认 `service.${control_name}`
	ReturnDataBool bool   // service 是否有返回数据
	DbConfig       string // 数据库填写

	LogOrSave string // 输出或者打印
}

type Req struct {
	ReqBool bool   // 是否自动创建绑定 3个以内推荐使用
	Req     string // 如果 有ParamReq 不填写 Req 获取 models.Req
}

// ${describe} 描述
// ${control_name} control 名称
// ${req} 请求名称
// ${log_or_save}
// ${data}

func GenerateController(c Control) string {
	var str = `
// ${describe}
func ${control_name}(c *gin.Context) {
	${req}
	if err := c.Bind(&req); err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 1, "msg": "request binding failed", "debug": err.Error()})
		return
	}

	${db_config}
	${service_str}
	if err != nil {
		${log_or_save}
		c.JSON(200, gin.H{"code": 2, "msg": "service operate failed", "debug": err.Error()})
		return
	}
	${out}
}
`
	str = strings.ReplaceAll(str, "${describe}", c.Describe)
	str = strings.ReplaceAll(str, "${req}", c.checkoutReq())
	str = strings.ReplaceAll(str, "${log_or_save}", c.checkoutLogOrSave())
	str = strings.ReplaceAll(str, "${service_str}", c.checkoutServiceStr())
	str = strings.ReplaceAll(str, "${out}", c.printOut())
	str = strings.ReplaceAll(str, "${db_config}", c.checkoutDbConfig())
	str = strings.ReplaceAll(str, "${control_name}", c.ControlName)

	return str
}

func (c Control) checkoutReq() string {
	if (Req{} == c.Req) {
		return ""
	}
	if c.Req.ReqBool {
		var str = `type Req struct {
	// TODO 请填写请求参数
	}
	var req Req`
		return str
	}
	return fmt.Sprintf("var req %+v", c.Req)
}

func (c Control) checkoutLogOrSave() string {
	if c.LogOrSave == `` {
		return `log.Printf("%+v",errors.WithStack(err))`
	}
	return c.LogOrSave
}

func (c Control) checkoutServiceStr() string {
	if c.ReturnDataBool {
		if c.ServiceStr != "" {
			return fmt.Sprintf(`data, err := %v(req, db)`, c.ServiceStr)
		}
		return fmt.Sprintf(`data, err := %v.%v(req, db)`, "service", c.ControlName)

	} else {

		if c.ServiceStr != "" {
			return fmt.Sprintf(`err := %v(req, db)`, c.ServiceStr)
		}
		return fmt.Sprintf(`err := %v.%v(req, db)`, "service", c.ControlName)
	}
}

func (c Control) printOut() string {
	if c.ReturnDataBool {
		return `
	log.Printf("way:%v ; req:%v ; data:%v ;", "${control_name}", req, data)
	c.JSON(200, gin.H{"code": 0, "msg": "success", "data": data})`
	}
	return `
	log.Printf("way:%v ; req:%v ; ", "${control_name}", req)
	c.JSON(200, gin.H{"code": 0, "msg": "success"})`
}

func (c Control) checkoutDbConfig() string {
	return fmt.Sprintf("db := %v", c.DbConfig)
}
