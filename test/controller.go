package test

import (
	"github.com/catbugdemo/auto_create/test/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

/*
	router.POST("/smb_vocher_list/add",handlers.HttpAddSmbVocherList)
	router.POST("/smb_vocher_list/get",handlers.HttpGetSmbVocherList)
	router.POST("/smb_vocher_list/list",,handlers.HttpListSmbVocherList)
	router.POST("/smb_vocher_list/update",handlers.HttpUpdateSmbVocherList)
	router.POST("/smb_vocher_list/delete",handlers.HttpDeleteSmbVocherList)
*/

// @Summary 卡券种类 list
// @title 后台接口
// @Tags 卡券
// @Router /smb_vocher_list/add [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.SmbVocherDetailLog true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpAddSmbVocherList(c *gin.Context) {
	var req model.SmbVocherDetailLog
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if err := req.Insert(c.MustGet("db").(*sqlx.DB)); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": req})
}

// @Summary 卡券种类 list
// @title 后台接口
// @Tags 卡券
// @Router /smb_vocher_list/get [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.SmbVocherDetailLog true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpGetSmbVocherList(c *gin.Context) {
	var req model.SmbVocherDetailLog
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if err := req.FirstById(c.MustGet("db").(*sqlx.DB), req.Id); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": req})
}

// @Summary 卡券种类 list
// @title 后台接口
// @Tags 卡券
// @Router /smb_vocher_list/list [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.SmbVocherDetailLog true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpListSmbVocherList(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")
	var req model.SmbVocherDetailLog
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	// generate limit
	list, count, err := req.FindByCount(c.MustGet("db").(*sqlx.DB), req.Condition(), limit, page)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": count})
}

// @Summary 卡券种类 list
// @title 后台接口
// @Tags 卡券
// @Router /smb_vocher_list/get [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.SmbVocherDetailLog true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpUpdateSmbVocherList(c *gin.Context) {
	var req model.SmbVocherDetailLog
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	id := req.Id
	req.Id = 0
	if err := req.UpdateById(c.MustGet("db").(*sqlx.DB), id); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": req})
}

// @Summary 卡券种类 list
// @title 后台接口
// @Tags 卡券
// @Router /smb_vocher_list/get [post]
//	@Security x-smb-jwt
//	@param x-smb-jwt header string true "Authorization"
// @param param body model.SmbVocherDetailLog true "用户请求参数"
// @Success 200 {object} JsonMsg
func HttpDelteSmbVocherList(c *gin.Context) {
	var req model.SmbVocherDetailLog
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if err := req.DeleteById(c.MustGet("db").(*sqlx.DB), req.Id); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": req})
}
