package auto

import "strings"

func AutoGenerateSqlxControl(c CRUD) string {

	c.initStruct()

	return c.formatSqlxCRUD()
}

func (s *St) formatSqlxCRUD() string {
	var str = `
/*
	router.POST("/${t_name}",handlers.HttpAdd${name})
	router.GET("/${t_name}/:id",handlers.HttpGet${name})
	router.POST("/${t_name}/list",handlers.HttpList${name})
	router.PATCH("/${t_name}/:id",handlers.HttpUpdate${name})
	router.DELETE("/${t_name}/:id",handlers.HttpDelete${name})
*/

// @Summary 新增
// @title 后台接口
// @Tags ${tag}
// @Router /${t_name} [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param req body models.${name} true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpAdd${name}(c *gin.Context) {
	var req models.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	if err := req.Insert(${db_config}); err != nil {
		c.JSON(200, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	c.JSON(200,gin.H{
		"error_no": CODE_SUCCESS_I,
		"error_msg": CODE_SUCCESS_S,
	})
}

// @Summary 根据 id 获取
// @title 后台接口
// @Tags ${tag}
// @Router /${t_name} [get]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body models.${name} true "用户请求参数"
// @param id path integer true "id"
// @Success 200 {object} JsonMsg{data=models.${name}}
func HttpGet${name}(c *gin.Context) {
	id := c.Param("id")
	var param models.${name}
	if err := param.FirstById(${db_config}, id); err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	c.JSON(200,gin.H{
		"data": param, 
		"error_no": CODE_SUCCESS_I,
		"error_msg": CODE_SUCCESS_S,
	})
}

// @Summary 获取 list
// @title 后台接口
// @Tags ${tag}
// @Router /${t_name}/list [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body models.${name} true "用户请求参数"
// @param size query string true "每页条数"
// @param page query string true "页数"
// @Success 200 {object} JsonMsg
func HttpList${name}(c *gin.Context) {
	size := c.Query("size")
	page := c.Query("page")
	var req models.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	count, err := req.Count(${db_config}, req.Condition())
	if err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	if count == 0 {
		c.JSON(200, gin.H{
			"data": map[string]interface{}{
				"list":  nil,
				"count": count,
			},
			"error_no": CODE_SUCCESS_I,
			"error_msg": CODE_SUCCESS_S,
		})
		return 
	}
	limit, offset := utils.ToLimitOffset(size, page, count)
	list, err := req.Find(c.MustGet(DB_CONFIG).(*sqlx.DB), req.Condition(), strconv.Itoa(limit), strconv.Itoa(offset))
	if err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
		"error_no": CODE_SUCCESS_I,
		"error_msg": CODE_SUCCESS_S,
	})
}

// @Summary 修改
// @title 后台接口
// @Tags ${tag}
// @Router /${t_name} [patch]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body models.${name} true "用户请求参数"
// @param id path integer true "id"
// @Success 200 {object} JsonMsg{data=models.${name}}
func HttpUpdate${name}(c *gin.Context) {
	id := c.Param("id")
	var req models.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	if err := req.UpdateById(${db_config}, id); err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	c.JSON(200,gin.H{
		"error_no": CODE_SUCCESS_I,
		"error_msg": CODE_SUCCESS_S,
	})
}

// @Summary 删除
// @title 后台接口
// @Tags 卡券
// @Router /${t_name} [delete]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body models.${name} true "用户请求参数"
// @param id path integer true "id"
// @Success 200 {object} JsonMsg
func HttpDelete${name}(c *gin.Context) {
	id := c.Param("id")
	var req models.${name}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	if err := req.DeleteById(${db_config}, id); err != nil {
		c.JSON(400, gin.H{
			"error_no": CODE_FAIL_I,
			"error_msg": fmt.Sprintf("%+v", err),
		})
		return
	}
	c.JSON(200,gin.H{
		"error_no": CODE_SUCCESS_I,
		"error_msg": CODE_SUCCESS_S,
	})
}
`
	str = strings.ReplaceAll(str, "${t_name}", s.parseLName())
	str = strings.ReplaceAll(str, "${name}", s.Info["name"].(string))
	str = strings.ReplaceAll(str, "${db_config}", s.DbConfig)
	str = strings.ReplaceAll(str, "${tag}", s.Info["tag"].(string))

	return str

}
