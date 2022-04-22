package auto

import "strings"

func AutoGenerateSqlxControl(c CRUD) string {

	c.initStruct()

	return c.formatSqlxCRUD()
}

func (s *St) formatSqlxCRUD() string {
	var str = `
/*
	router.POST("/${t_name}/add",handlers.HttpAdd${name})
	router.POST("/${t_name}/get",handlers.HttpGet${name})
	router.POST("/${t_name}/list",handlers.HttpList${name})
	router.POST("/${t_name}/update",handlers.HttpUpdate${name})
	router.POST("/${t_name}/delete",handlers.HttpDelete${name})
*/

// @Summary 新增
// @title 后台接口
// @Tags 卡券
// @Router /${t_name}/add [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.${name} true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpAdd${name}(c *gin.Context) {
	var req model.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		${bind_json}
		return
	}
	if err := req.Insert(${db_config}); err != nil {
		${service_json}
		return
	}
	${success_json}
}

// @Summary 获取
// @title 后台接口
// @Tags 卡券
// @Router /${t_name}/get [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.${name} true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpGet${name}(c *gin.Context) {
	var req model.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		${bind_json}
		return
	}
	if err := req.FirstById(${db_config}, req.Id); err != nil {
		${service_json}
		return
	}
	${success_json}
}

// @Summary 获取 list
// @title 后台接口
// @Tags 卡券
// @Router /${t_name}/list [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.${name} true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpList${name}(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")
	var req model.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		${bind_json}
		return
	}
	list, count, err := req.FindByCount(${db_config}, req.Condition())
	if err != nil {
		${service_json}
		return
	}
	${success_json}
}

// @Summary 修改
// @title 后台接口
// @Tags 卡券
// @Router /${t_name}/update [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.${name} true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpUpdate${name}(c *gin.Context) {
	var req model.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		${bind_json}
		return
	}
	id := req.Id
	req.Id = 0
	if err := req.UpdateById(${db_config}, id); err != nil {
		${service_json}
		return
	}
	${success_json}
}

// @Summary 删除
// @title 后台接口
// @Tags 卡券
// @Router /${t_name}/delete [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.${name} true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpDelete${name}(c *gin.Context) {
	var req model.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		${bind_json}
		return
	}
	if err := req.DeleteById(${db_config}, req.Id); err != nil {
		${service_json}
		return
	}
	${success_json}
}
`
	str = strings.ReplaceAll(str, "${t_name}", s.parseLName())
	str = strings.ReplaceAll(str, "${name}", s.Info["name"].(string))
	str = strings.ReplaceAll(str, "${db_config}", s.DbConfig)

	str = strings.ReplaceAll(str, "${bind_json}", s.Info["bind_json"].(string))
	str = strings.ReplaceAll(str, "${service_json}", s.Info["service_json"].(string))
	str = strings.ReplaceAll(str, "${success_json}", s.Info["success_json"].(string))
	return str

}
