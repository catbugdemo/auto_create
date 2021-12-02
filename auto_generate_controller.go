package auto

/*
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
	Req         string `json:"req"`       // 请求
	DbConfig    string `json:"db_config"` // 数据库配置

	LogBind    string `json:"log_bind"`    // 绑定错误返回
	LogService string `json:"log_service"` // service 层错误返回
	LogReturn  string `json:"log_return"`  // return成功返回
}

type Swag struct {
	Describe string `json:"describe"` // 描述
	Tags     string `json:"tags"`     // 标签
	Router   string `json:"router"`   // 路由 默认 post
	Security string `json:"security"` // 是否需要 header_token
}

// ${describe} 描述
// ${tags} 分类
// ${router}
// ${security}

// ${control_name} control 名称
// ${req} 请求名称

func GenerateController(o Model) string {
	var str = `
// @Summary ${describe}
// @title 后台接口
// @Tags ${tags}
// @Router ${router} [post] ${security}
// @param param body 请写req true "用户请求参数"
// @Success 200 {object} JsonMsg
func ${control_name}(c *gin.Context) {
	${req}
	data, err := service.${control_name}(req,${db_config}})
	if err != nil {
		${log_service}
		return
	}
	${log_return}
}
`
	str = strings.ReplaceAll(str, "${describe}", o.Swag.Describe)
	str = strings.ReplaceAll(str, "${tags}", o.Swag.Tags)
	str = strings.ReplaceAll(str, "${router}", o.Swag.Router)
	str = strings.ReplaceAll(str, "${security}", o.Swag.DealSecurity())

	str = strings.ReplaceAll(str, "${req}", o.Control.checkoutReq())
	str = strings.ReplaceAll(str, "${control_name}", o.ControlName)
	str = strings.ReplaceAll(str, "${db_config}", o.Control.checkoutDbConfig())

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

func (c Control) checkoutReq() string {
	if c.Req == "" {
		return ""
	}
	return `var req
	if err := c.Bind(&req); err != nil {
		${log_bind}
		return
	}`
}

func (c Control) checkoutDbConfig() string {
	return fmt.Sprintf("db := %v", c.DbConfig)
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
*/
