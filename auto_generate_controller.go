package auto

import (
	"fmt"
	"strings"
)

type Model struct {
	Swag
	Control
}

type Control struct {
	ControlName string `json:"name"`      // controller 名
	DbConfig    string `json:"db_config"` // 数据库配置

	LogBind    string `json:"log_bind"`    // 绑定错误返回
	LogService string `json:"log_service"` // service 层错误返回
	LogReturn  string `json:"log_return"`  // return成功返回
}

type Swag struct {
	Security string `json:"security"` // 是否需要 header_token
}

// ${router}
// ${security}
// ${control_name} control 名称
// ${req} 请求名称

func GenerateController(o Model) string {
	var str = `
type ${control_name}Req struct{
	
}

// @Summary 描述
// @title 后台接口
// @Tags 标签分类
// @Router /test [post] ${security}
// @param param body ${control_name}Req true "用户请求参数"
// @Success 200 {object} JsonMsg
func ${control_name}(c *gin.Context) {
	timeNow := time.Now()
	var req ${control_name}Req
	if err := c.BindJSON(&req); err != nil {
		${log_bind}
		return
	}

	data, err := service.${control_name}(req,${db_config})
	if err != nil {
		${log_service}
		return
	}
	${log_return}
	return
}
`
	str = strings.ReplaceAll(str, "${security}", o.Swag.DealSecurity())
	str = strings.ReplaceAll(str, "${control_name}", o.ControlName)
	str = strings.ReplaceAll(str, "${db_config}", o.Control.checkoutDbConfig())
	str = strings.ReplaceAll(str, "${log_bind}", o.LogBind)
	str = strings.ReplaceAll(str, "${log_service}", o.LogService)
	str = strings.ReplaceAll(str, "${log_return}", o.LogReturn)
	return str
}

func (o Swag) DealSecurity() string {
	if o.Security == "" {
		return ""
	}
	var str = `
//	@Security ${security}
//	@param ${security} header string true "Authorization"`
	return strings.ReplaceAll(str, "${security}", o.Security)
}

func (c Control) checkoutDbConfig() string {
	return fmt.Sprintf("db := %v", c.DbConfig)
}
