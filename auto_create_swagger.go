package auto

import "strings"

// ${name} 该结构体中文名
// ${tags} 数据库名称
// ${router} 路由名称
// ${req} 请求结构体

type SwaggerInit struct {
	Name   string `json:"name"`
	Tags   string `json:"tags"`
	Router string `json:"router"`
	Req    string `json:"req"`
}

type SwaggerNormal struct {
	Describe      string `json:"describe"`
	TableName     string `json:"table_name"`
	IfHeaderToken bool   `json:"if_header_token"`
	Req           string `json:"if_req"`
	RouterUrl     string `json:"router_url"`
}

func GenerateNormalSwagger(swag SwaggerNormal) string {
	var str = `
// @Summary ${describe}
// @title 腾讯插件接口
// @Tags ${table_name}
// @Router ${router} [post] ${if_header_token}
// @param param body ${req} true "用户请求参数"
// @Success 200 {object} JsonMsg
`
	str = strings.ReplaceAll(str, "${describe}", swag.Describe)
	str = strings.ReplaceAll(str, "${table_name}", swag.TableName)
	str = strings.ReplaceAll(str, "${req}", swag.Req)
	str = strings.ReplaceAll(str, "${if_header_token}", IfHeaderToken(swag.IfHeaderToken))
	str = strings.ReplaceAll(str, "${router}", swag.RouterUrl)
	return str
}

func IfHeaderToken(check bool) string {
	if check {
		return `
//	@Security x-ytf-jwt
//	@param x-ytf-jwt header string true "Authorization"`
	}
	return ""
}

func GenerateCURDSwagger(swag SwaggerInit) string {
	var str = `
// @Summary 新增${name}
// @title 腾讯插件接口
// @Tags ${tags}
// @Router /${router} [post]
//	@Security x-ytf-jwt
//	@param x-ytf-jwt header string true "Authorization"
// @param param body ${req} true "用户请求参数"
// @Success 200 {object} JsonMsg





// @Summary 获取${name}
// @title 腾讯插件接口
// @Tags ${tags}
// @Router /${router}/{id} [get]
//	@Security x-ytf-jwt
//	@param x-ytf-jwt header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} JsonMsg




// @Summary 修改${name}
// @title 腾讯插件接口
// @Tags ${tags}
// @Router /${router}/{id} [patch]
//	@Security x-ytf-jwt
//	@param x-ytf-jwt header string true "Authorization"
// @param id path string true "id"
// @param param body ${req} true "用户请求参数"
// @Success 200 {object} JsonMsg



// @Summary 删除${name}
// @title 腾讯插件接口
// @Tags ${tags}
// @Router /${router}/{id} [delete]
//	@Security x-ytf-jwt
//	@param x-ytf-jwt header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} JsonMsg




// @Summary 获取${name}list
// @title 腾讯插件接口
// @Tags ${tags}
// @Router /${router} [get]
//	@Security x-ytf-jwt
//	@param x-ytf-jwt header string true "Authorization"
// @Success 200 {object} JsonMsgList

`
	str = strings.ReplaceAll(str, "${name}", swag.Name)
	str = strings.ReplaceAll(str, "${tags}", swag.Tags)
	str = strings.ReplaceAll(str, "${router}", swag.Router)
	str = strings.ReplaceAll(str, "${req}", swag.Req)
	return str
}
